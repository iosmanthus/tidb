package domain

import (
	"os"

	"github.com/pingcap/tidb/config"
)

const EnvVarKeyspaceName = "KEYSPACE_NAME"

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
