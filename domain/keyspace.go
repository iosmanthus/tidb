package domain

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/kv"
)

const (

	// EnvVarKeyspaceName is the system env name for keyspace name.
	EnvVarKeyspaceName = "KEYSPACE_NAME"

	// tidbKeyspaceEtcdPathPrefix is the keyspace prefix for etcd namespace
	tidbKeyspaceEtcdPathPrefix = "/keyspaces/tidb/"
)

func GetKeyspaceNameBySettings() (keyspaceName string) {

	// The cfg.keyspaceName get higher weights than KEYSPACE_NAME in system env.
	keyspaceName = config.GetGlobalConfig().KeyspaceName
	if !IsKeyspaceNameEmpty(keyspaceName) {
		return keyspaceName
	}

	keyspaceName = os.Getenv(EnvVarKeyspaceName)
	return keyspaceName
}

func IsKeyspaceNameEmpty(keyspaceName string) bool {
	return keyspaceName == ""
}

// GetKeyspacePathPrefix return the keyspace prefix path for etcd namespace
func GetKeyspacePathPrefix(keyspaceId uint32) string {
	path := fmt.Sprintf(tidbKeyspaceEtcdPathPrefix+"%d"+"/", keyspaceId)
	return path
}

// KeyspaceIdBytesToUint32 is used to convert byte array to uint32
func KeyspaceIdBytesToUint32(b []byte) uint32 {
	c := make([]byte, 4)
	copy(c[1:4], b[0:3])
	return binary.BigEndian.Uint32(c)
}

func isKvStorageKeyspaceSet(store kv.Storage) bool {
	return store.GetCodec().GetKeyspace() != nil
}
