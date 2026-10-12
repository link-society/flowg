package config_test

import (
	"testing"

	"os"
	"path/filepath"
	"strings"

	"github.com/stretchr/testify/require"

	"link-society.com/flowg/internal/app/config"
)

// clearEnv blanks the FLOWG_* environment variables inherited from the shell,
// so that tests only see the variables they set. An empty variable is treated
// as unset by the configuration.
func clearEnv(t *testing.T) {
	t.Helper()

	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "FLOWG_") {
			t.Setenv(key, "")
		}
	}
}

func loadConfig(t *testing.T, source string) *config.RootConfig {
	t.Helper()
	clearEnv(t)

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
	clearEnv(t)

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

func TestDefaultConfig_ReadsEnvAtCallTime(t *testing.T) {
	clearEnv(t)
	t.Setenv("FLOWG_HTTP_BIND_ADDRESS", ":1234")
	t.Setenv("FLOWG_BADGER_LOG_DIR", "/tmp/flowg/logs")

	cfg := config.DefaultConfig()

	require.Equal(t, ":1234", cfg.Services.Http.BindAddress)
	require.Equal(t, "/tmp/flowg/logs", cfg.Storage.Backend.BadgerDB.LogDir)
}

func TestValidate_EnvOnlyConfig(t *testing.T) {
	testCases := []struct {
		name string
		key  string
		val  string
		want string
	}{
		{
			name: "invalid syslog protocol",
			key:  "FLOWG_SYSLOG_PROTOCOL",
			val:  "bogus",
			want: "invalid syslog protocol: bogus",
		},
		{
			name: "unsupported storage backend",
			key:  "FLOWG_STORAGE_BACKEND",
			val:  "bogus",
			want: `unsupported storage backend: "bogus"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearEnv(t)
			t.Setenv(tc.key, tc.val)

			cfg := config.DefaultConfig()
			require.EqualError(t, cfg.Validate(), tc.want)
		})
	}
}
