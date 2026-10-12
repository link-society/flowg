package cmd_test

import (
	"testing"

	"os"
	"path/filepath"

	"github.com/stretchr/testify/require"

	"link-society.com/flowg/cmd/flowg-server/cmd"
	"link-society.com/flowg/internal/app/server"
)

func loadConfig(t *testing.T, source string) *cmd.RootConfig {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.hcl")
	require.NoError(t, os.WriteFile(path, []byte(source), 0600))

	config := cmd.DefaultConfig()
	require.NoError(t, config.Load(path))

	return config
}

func TestLoad_StorageBackendBadgerDbFromFile(t *testing.T) {
	config := loadConfig(t, `
storage {
  backend "badgerdb" {
    auth_dir   = "/tmp/flowg/auth"
    config_dir = "/tmp/flowg/config"
    log_dir    = "/tmp/flowg/logs"
  }
}
`)

	opts, err := config.AsServerOptions()
	require.NoError(t, err)
	require.Equal(
		t,
		&server.BadgerDbStorageOptions{
			AuthDir:   "/tmp/flowg/auth",
			ConfigDir: "/tmp/flowg/config",
			LogDir:    "/tmp/flowg/logs",
		},
		opts.Storage,
	)
}

func TestLoad_StorageBackendFoundationDbFromFile(t *testing.T) {
	config := loadConfig(t, `
storage {
  backend "foundationdb" {
    cluster_file      = "/tmp/flowg/fdb.cluster"
    connection_string = "flowg:test@127.0.0.1:4500"
    key_space         = "flowg-test"
  }
}
`)

	opts, err := config.AsServerOptions()
	require.NoError(t, err)
	require.Equal(
		t,
		&server.FoundationDbStorageOptions{
			ClusterFile:      "/tmp/flowg/fdb.cluster",
			ConnectionString: "flowg:test@127.0.0.1:4500",
			KeySpace:         "flowg-test",
		},
		opts.Storage,
	)
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

	config := cmd.DefaultConfig()
	require.ErrorContains(
		t,
		config.Load(path),
		"FoundationDB backend requires either a cluster file or a connection string",
	)
}
