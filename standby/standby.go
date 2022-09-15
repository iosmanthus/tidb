package standby

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pingcap/tidb/util/logutil"
	"go.uber.org/zap"
)

const (
	standbyState   = "standby"
	activatedState = "activated"
)

var (
	mu           sync.RWMutex
	state        = standbyState
	keyspaceName string

	// activationTimeout specifies the maximum allowed time for tidb to activate from standby mode.
	activationTimeout uint
)

var (
	activateCh      = make(chan struct{}, 1)
	serverIsReadyCh = make(chan struct{})
)

// StandbyHandler returns a handler to query tidb pool status or activate or exit the tidb server.
func StandbyHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/tidb-pool/status", statusHandler)
	mux.HandleFunc("/tidb-pool/activate", func(w http.ResponseWriter, r *http.Request) {
		type activateRequest struct {
			KeyspaceName string `json:"keyspace_name"`
		}
		var req activateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.KeyspaceName == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mu.Lock()
		if state == standbyState {
			state = activatedState
			keyspaceName = req.KeyspaceName
			activateCh <- struct{}{}
		} else if keyspaceName != req.KeyspaceName {
			mu.Unlock()
			w.WriteHeader(http.StatusPreconditionFailed)
			w.Write([]byte("server is not in standby mode"))
			return
		}
		// if client tries to activate with same keyspace name, wait for ready signal and return 200.
		mu.Unlock()

		// If no limit posted on activation time, wait for serverIsReady indefinitely.
		if activationTimeout == 0 {
			<-serverIsReadyCh
			statusHandler(w, r)
			return
		}
		select {
		case <-time.After(time.Duration(activationTimeout) * time.Second):
			logutil.BgLogger().Warn("timeout waiting for activation")
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("timeout waiting for activation"))
			os.Exit(1)
		case <-serverIsReadyCh:
			statusHandler(w, r)
		}
	})
	mux.HandleFunc("/tidb-pool/exit", func(w http.ResponseWriter, r *http.Request) {
		logutil.BgLogger().Info("receiving exit signal, exiting...")
		os.Exit(0)
	})
	return mux
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"state": "%s", "keyspace_name": "%s"}`, state, keyspaceName)
}

var server *http.Server

// StartStandby starts a http server to listen and wait for activation signal.
func StartStandby(host string, port uint, timeout uint) string {
	mux := StandbyHandler()
	// handle liveness probe.
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
	}
	activationTimeout = timeout
	logutil.BgLogger().Info("tidb-server is now running as standby, waiting for activation...", zap.String("addr", server.Addr))
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logutil.BgLogger().Warn("failed to start tidb-server as standby", zap.Error(err))
			os.Exit(1)
		}
	}()

	<-activateCh

	mu.RLock()
	defer mu.RUnlock()
	return keyspaceName
}

// EndStandby is used to notify the temp http server that the tidb server is ready.
func EndStandby() error {
	close(serverIsReadyCh)
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}
	return nil
}
