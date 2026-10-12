package config_test

import (
	"testing"

	"os"
	"path/filepath"

	"github.com/stretchr/testify/require"

	"link-society.com/flowg/internal/app/config"
)

func loadConfig(t *testing.T, source string) *config.RootConfig {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.hcl")
	require.NoError(t, os.WriteFile(path, []byte(source), 0600))

	cfg := config.DefaultConfig()
	require.NoError(t, cfg.Load(path))

	return cfg
}

func TestLoad_StorageBackendBadgerDbFromFile(t *testing.T) {
	cfg := loadConfig(t, `
storage {
  backend "badgerdb" {
    auth_dir   = "/tmp/flowg/auth"
    config_dir = "/tmp/flowg/config"
    log_dir    = "/tmp/flowg/logs"
  }
}
`)

	require.Equal(
		t,
		&config.StorageBackendBadgerDbConfig{
			AuthDir:   "/tmp/flowg/auth",
			ConfigDir: "/tmp/flowg/config",
			LogDir:    "/tmp/flowg/logs",
		},
		cfg.Storage.Backend.BadgerDB,
	)
	require.Nil(t, cfg.Storage.Backend.FoundationDB)
}

func TestLoad_StorageBackendFoundationDbFromFile(t *testing.T) {
	cfg := loadConfig(t, `
storage {
  backend "foundationdb" {
    cluster_file      = "/tmp/flowg/fdb.cluster"
    connection_string = "flowg:test@127.0.0.1:4500"
    key_space         = "flowg-test"
  }
}
`)

	require.Equal(
		t,
		&config.StorageBackendFoundationDbConfig{
			ClusterFile:      "/tmp/flowg/fdb.cluster",
			ConnectionString: "flowg:test@127.0.0.1:4500",
			KeySpace:         "flowg-test",
		},
		cfg.Storage.Backend.FoundationDB,
	)
	require.Nil(t, cfg.Storage.Backend.BadgerDB)
}

func TestLoad_StorageBackendValidationError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.hcl")
	require.NoError(t, os.WriteFile(path, []byte(`
storage {
  backend "foundationdb" {
    cluster_file      = ""
    connection_string = ""
  }
}
`), 0600))

	cfg := config.DefaultConfig()
	require.ErrorContains(
		t,
		cfg.Load(path),
		"FoundationDB backend requires either a cluster file or a connection string",
	)
}
