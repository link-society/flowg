package config_test

import (
	"testing"

	"os"
	"path/filepath"

	"github.com/stretchr/testify/require"

	"link-society.com/flowg/internal/app/config"
)

// block formats an HCL block with the given header and body.
func block(header string, body ...string) string {
	content := ""
	for _, line := range body {
		content += line + "\n"
	}

	return header + " {\n" + content + "}"
}

// loadConfigWithEnv loads the default configuration, overridden by the given
// configuration file (skipped when empty) and environment variables.
func loadConfigWithEnv(
	t *testing.T,
	source string,
	env map[string]string,
) *config.RootConfig {
	t.Helper()
	clearEnv(t)

	for key, value := range env {
		t.Setenv(key, value)
	}

	path := ""
	if source != "" {
		path = filepath.Join(t.TempDir(), "config.hcl")
		require.NoError(t, os.WriteFile(path, []byte(source), 0600))
	}

	cfg := config.DefaultConfig()
	require.NoError(t, cfg.Load(path))

	return cfg
}

func TestLoad_Precedence(t *testing.T) {
	seedAuth := func(lines ...string) string {
		return block("storage", block("seed", block("auth", lines...)))
	}

	settings := []struct {
		name string
		// baseEnv is applied in every scenario.
		baseEnv map[string]string
		// file is a configuration file setting the value to fromFile.
		file string
		// env is the environment variable setting the value to fromEnv.
		env      string
		envValue string

		get       func(*config.RootConfig) any
		byDefault any
		fromFile  any
		fromEnv   any
	}{
		{
			name:      "demo",
			file:      `demo = true`,
			env:       "FLOWG_DEMO_MODE",
			envValue:  "false",
			get:       func(c *config.RootConfig) any { return c.DemoMode },
			byDefault: false,
			fromFile:  true,
			fromEnv:   false,
		},
		{
			name:      "logging.level",
			file:      block("logging", `level = "warn"`),
			env:       "FLOWG_LOGLEVEL",
			envValue:  "error",
			get:       func(c *config.RootConfig) any { return c.Logging.Level },
			byDefault: "info",
			fromFile:  "warn",
			fromEnv:   "error",
		},
		{
			name:      "logging.verbose",
			file:      block("logging", `verbose = true`),
			env:       "FLOWG_VERBOSE",
			envValue:  "false",
			get:       func(c *config.RootConfig) any { return c.Logging.Verbose },
			byDefault: false,
			fromFile:  true,
			fromEnv:   false,
		},
		{
			name:      "services.http.bind",
			file:      block("services", block("http", `bind = ":1001"`)),
			env:       "FLOWG_HTTP_BIND_ADDRESS",
			envValue:  ":1002",
			get:       func(c *config.RootConfig) any { return c.Services.Http.BindAddress },
			byDefault: ":5080",
			fromFile:  ":1001",
			fromEnv:   ":1002",
		},
		{
			name:      "services.http.mount",
			file:      block("services", block("http", `mount = "/file"`)),
			env:       "FLOWG_HTTP_MOUNT_PATH",
			envValue:  "/env",
			get:       func(c *config.RootConfig) any { return c.Services.Http.MountPath },
			byDefault: "/",
			fromFile:  "/file",
			fromEnv:   "/env",
		},
		{
			name:      "services.management.bind",
			file:      block("services", block("management", `bind = ":1001"`)),
			env:       "FLOWG_MGMT_BIND_ADDRESS",
			envValue:  ":1002",
			get:       func(c *config.RootConfig) any { return c.Services.Management.BindAddress },
			byDefault: ":9113",
			fromFile:  ":1001",
			fromEnv:   ":1002",
		},
		{
			name:      "services.syslog.bind",
			file:      block("services", block("syslog", `bind = ":1001"`)),
			env:       "FLOWG_SYSLOG_BIND_ADDRESS",
			envValue:  ":1002",
			get:       func(c *config.RootConfig) any { return c.Services.Syslog.Bind },
			byDefault: ":5514",
			fromFile:  ":1001",
			fromEnv:   ":1002",
		},
		{
			name:      "services.syslog.protocol",
			file:      block("services", block("syslog", `protocol = "tcp"`)),
			env:       "FLOWG_SYSLOG_PROTOCOL",
			envValue:  "udp",
			get:       func(c *config.RootConfig) any { return c.Services.Syslog.Protocol },
			byDefault: "udp",
			fromFile:  "tcp",
			fromEnv:   "udp",
		},
		{
			name:      "services.syslog.initial_allowed_origins",
			file:      block("services", block("syslog", `initial_allowed_origins = ["10.0.0.1"]`)),
			env:       "FLOWG_SYSLOG_INITIAL_ALLOWED_ORIGINS",
			envValue:  "10.0.0.2, 10.0.0.3",
			get:       func(c *config.RootConfig) any { return c.Services.Syslog.InitialAllowedOrigins },
			byDefault: []string{},
			fromFile:  []string{"10.0.0.1"},
			fromEnv:   []string{"10.0.0.2", "10.0.0.3"},
		},
		{
			name:      "storage.backend.badgerdb.auth_dir",
			file:      block("storage", block(`backend "badgerdb"`, `auth_dir = "/file/auth"`)),
			env:       "FLOWG_BADGER_AUTH_DIR",
			envValue:  "/env/auth",
			get:       func(c *config.RootConfig) any { return c.Storage.Backend.BadgerDB.AuthDir },
			byDefault: "./data/auth",
			fromFile:  "/file/auth",
			fromEnv:   "/env/auth",
		},
		{
			name:      "storage.backend.badgerdb.log_dir",
			file:      block("storage", block(`backend "badgerdb"`, `log_dir = "/file/logs"`)),
			env:       "FLOWG_BADGER_LOG_DIR",
			envValue:  "/env/logs",
			get:       func(c *config.RootConfig) any { return c.Storage.Backend.BadgerDB.LogDir },
			byDefault: "./data/logs",
			fromFile:  "/file/logs",
			fromEnv:   "/env/logs",
		},
		{
			name:      "storage.backend.badgerdb.config_dir",
			file:      block("storage", block(`backend "badgerdb"`, `config_dir = "/file/config"`)),
			env:       "FLOWG_BADGER_CONFIG_DIR",
			envValue:  "/env/config",
			get:       func(c *config.RootConfig) any { return c.Storage.Backend.BadgerDB.ConfigDir },
			byDefault: "./data/config",
			fromFile:  "/file/config",
			fromEnv:   "/env/config",
		},
		{
			name:      "storage.backend.foundationdb.cluster_file",
			baseEnv:   map[string]string{"FLOWG_STORAGE_BACKEND": "foundationdb"},
			file:      block("storage", block(`backend "foundationdb"`, `cluster_file = "/file/fdb.cluster"`)),
			env:       "FLOWG_FOUNDATIONDB_CLUSTER_FILE",
			envValue:  "/env/fdb.cluster",
			get:       func(c *config.RootConfig) any { return c.Storage.Backend.FoundationDB.ClusterFile },
			byDefault: "/etc/foundationdb/fdb.cluster",
			fromFile:  "/file/fdb.cluster",
			fromEnv:   "/env/fdb.cluster",
		},
		{
			name:      "storage.backend.foundationdb.connection_string",
			baseEnv:   map[string]string{"FLOWG_STORAGE_BACKEND": "foundationdb"},
			file:      block("storage", block(`backend "foundationdb"`, `connection_string = "file:file@127.0.0.1:4500"`)),
			env:       "FLOWG_FOUNDATIONDB_CONNECTION_STRING",
			envValue:  "env:env@127.0.0.1:4500",
			get:       func(c *config.RootConfig) any { return c.Storage.Backend.FoundationDB.ConnectionString },
			byDefault: "",
			fromFile:  "file:file@127.0.0.1:4500",
			fromEnv:   "env:env@127.0.0.1:4500",
		},
		{
			name:      "storage.backend.foundationdb.key_space",
			baseEnv:   map[string]string{"FLOWG_STORAGE_BACKEND": "foundationdb"},
			file:      block("storage", block(`backend "foundationdb"`, `key_space = "file"`)),
			env:       "FLOWG_FOUNDATIONDB_KEY_SPACE",
			envValue:  "env",
			get:       func(c *config.RootConfig) any { return c.Storage.Backend.FoundationDB.KeySpace },
			byDefault: "flowg",
			fromFile:  "file",
			fromEnv:   "env",
		},
		{
			name:      "storage.seed.auth.initial_user",
			file:      seedAuth(`initial_user = "file"`, `initial_password = "root"`),
			env:       "FLOWG_AUTH_INITIAL_USER",
			envValue:  "env",
			get:       func(c *config.RootConfig) any { return c.Storage.Seed.Auth.InitialUser },
			byDefault: "root",
			fromFile:  "file",
			fromEnv:   "env",
		},
		{
			name:      "storage.seed.auth.initial_password",
			file:      seedAuth(`initial_user = "root"`, `initial_password = "file"`),
			env:       "FLOWG_AUTH_INITIAL_PASSWORD",
			envValue:  "env",
			get:       func(c *config.RootConfig) any { return c.Storage.Seed.Auth.InitialPassword },
			byDefault: "root",
			fromFile:  "file",
			fromEnv:   "env",
		},
		{
			name:      "storage.seed.auth.reset_user",
			file:      seedAuth(`initial_user = "root"`, `initial_password = "root"`, `reset_user = "file"`),
			env:       "FLOWG_AUTH_RESET_USER",
			envValue:  "env",
			get:       func(c *config.RootConfig) any { return c.Storage.Seed.Auth.ResetUser },
			byDefault: "",
			fromFile:  "file",
			fromEnv:   "env",
		},
		{
			name:      "storage.seed.auth.reset_password",
			file:      seedAuth(`initial_user = "root"`, `initial_password = "root"`, `reset_password = "file"`),
			env:       "FLOWG_AUTH_RESET_PASSWORD",
			envValue:  "env",
			get:       func(c *config.RootConfig) any { return c.Storage.Seed.Auth.ResetPassword },
			byDefault: "",
			fromFile:  "file",
			fromEnv:   "env",
		},
	}

	for _, s := range settings {
		withEnv := map[string]string{s.env: s.envValue}
		for key, value := range s.baseEnv {
			withEnv[key] = value
		}

		scenarios := []struct {
			name string
			file string
			env  map[string]string
			want any
		}{
			{name: "default", file: "", env: s.baseEnv, want: s.byDefault},
			{name: "file over default", file: s.file, env: s.baseEnv, want: s.fromFile},
			{name: "env over default", file: "", env: withEnv, want: s.fromEnv},
			{name: "env over file", file: s.file, env: withEnv, want: s.fromEnv},
		}

		for _, sc := range scenarios {
			t.Run(s.name+"/"+sc.name, func(t *testing.T) {
				cfg := loadConfigWithEnv(t, sc.file, sc.env)
				require.Equal(t, sc.want, s.get(cfg))
			})
		}
	}
}

func TestLoad_TlsPrecedence(t *testing.T) {
	tls := func(service string, lines ...string) string {
		return block("services", block(service, block("tls", lines...)))
	}

	testCases := []struct {
		name string
		file string
		env  map[string]string
		get  func(*config.RootConfig) any
		want any
	}{
		{
			name: "http disabled by default",
			get:  func(c *config.RootConfig) any { return c.Services.Http.Tls },
			want: (*config.HttpTlsConfig)(nil),
		},
		{
			name: "http enabled by env",
			env: map[string]string{
				"FLOWG_HTTP_TLS_ENABLED": "true",
				"FLOWG_HTTP_TLS_CERT":    "env.crt",
				"FLOWG_HTTP_TLS_KEY":     "env.key",
			},
			get:  func(c *config.RootConfig) any { return c.Services.Http.Tls },
			want: &config.HttpTlsConfig{Cert: "env.crt", Key: "env.key"},
		},
		{
			name: "http enabled by file, overridden by env",
			file: tls("http", `cert = "file.crt"`, `key = "file.key"`),
			env:  map[string]string{"FLOWG_HTTP_TLS_CERT": "env.crt"},
			get:  func(c *config.RootConfig) any { return c.Services.Http.Tls },
			want: &config.HttpTlsConfig{Cert: "env.crt", Key: "file.key"},
		},
		{
			name: "http enabled by file, disabled by env",
			file: tls("http", `cert = "file.crt"`, `key = "file.key"`),
			env:  map[string]string{"FLOWG_HTTP_TLS_ENABLED": "false"},
			get:  func(c *config.RootConfig) any { return c.Services.Http.Tls },
			want: (*config.HttpTlsConfig)(nil),
		},
		{
			name: "management enabled by env",
			env: map[string]string{
				"FLOWG_MGMT_TLS_ENABLED": "true",
				"FLOWG_MGMT_TLS_CERT":    "env.crt",
				"FLOWG_MGMT_TLS_KEY":     "env.key",
			},
			get:  func(c *config.RootConfig) any { return c.Services.Management.Tls },
			want: &config.HttpTlsConfig{Cert: "env.crt", Key: "env.key"},
		},
		{
			name: "management enabled by file, overridden by env",
			file: tls("management", `cert = "file.crt"`, `key = "file.key"`),
			env:  map[string]string{"FLOWG_MGMT_TLS_KEY": "env.key"},
			get:  func(c *config.RootConfig) any { return c.Services.Management.Tls },
			want: &config.HttpTlsConfig{Cert: "file.crt", Key: "env.key"},
		},
		{
			name: "management enabled by file, disabled by env",
			file: tls("management", `cert = "file.crt"`, `key = "file.key"`),
			env:  map[string]string{"FLOWG_MGMT_TLS_ENABLED": "false"},
			get:  func(c *config.RootConfig) any { return c.Services.Management.Tls },
			want: (*config.HttpTlsConfig)(nil),
		},
		{
			name: "syslog enabled by env",
			env: map[string]string{
				"FLOWG_SYSLOG_PROTOCOL":    "tcp",
				"FLOWG_SYSLOG_TLS_ENABLED": "true",
				"FLOWG_SYSLOG_TLS_CERT":    "env.crt",
				"FLOWG_SYSLOG_TLS_KEY":     "env.key",
				"FLOWG_SYSLOG_TLS_AUTH":    "true",
			},
			get:  func(c *config.RootConfig) any { return c.Services.Syslog.Tls },
			want: &config.SyslogTlsConfig{Cert: "env.crt", Key: "env.key", AuthEnabled: true},
		},
		{
			name: "syslog enabled by file, overridden by env",
			file: block("services", block("syslog",
				`protocol = "tcp"`,
				block("tls", `cert = "file.crt"`, `key = "file.key"`, `auth = true`),
			)),
			env:  map[string]string{"FLOWG_SYSLOG_TLS_AUTH": "false"},
			get:  func(c *config.RootConfig) any { return c.Services.Syslog.Tls },
			want: &config.SyslogTlsConfig{Cert: "file.crt", Key: "file.key", AuthEnabled: false},
		},
		{
			name: "syslog enabled by file, disabled by env",
			file: block("services", block("syslog",
				`protocol = "tcp"`,
				block("tls", `cert = "file.crt"`, `key = "file.key"`),
			)),
			env:  map[string]string{"FLOWG_SYSLOG_TLS_ENABLED": "false"},
			get:  func(c *config.RootConfig) any { return c.Services.Syslog.Tls },
			want: (*config.SyslogTlsConfig)(nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := loadConfigWithEnv(t, tc.file, tc.env)
			require.Equal(t, tc.want, tc.get(cfg))
		})
	}
}

func TestLoad_StorageBackendPrecedence(t *testing.T) {
	t.Run("env selects foundationdb", func(t *testing.T) {
		cfg := loadConfigWithEnv(t, "", map[string]string{
			"FLOWG_STORAGE_BACKEND":        "foundationdb",
			"FLOWG_FOUNDATIONDB_KEY_SPACE": "env",
		})

		require.Equal(t, config.StorageBackendFoundationDb, cfg.Storage.Backend.Backend)
		require.Equal(
			t,
			&config.StorageBackendFoundationDbConfig{
				ClusterFile: "/etc/foundationdb/fdb.cluster",
				KeySpace:    "env",
			},
			cfg.Storage.Backend.FoundationDB,
		)
		require.Nil(t, cfg.Storage.Backend.BadgerDB)
	})

	t.Run("env overrides the backend selected by file", func(t *testing.T) {
		cfg := loadConfigWithEnv(
			t,
			block("storage", block(`backend "foundationdb"`, `cluster_file = "/file/fdb.cluster"`)),
			map[string]string{
				"FLOWG_STORAGE_BACKEND": "badgerdb",
				"FLOWG_BADGER_LOG_DIR":  "/env/logs",
			},
		)

		require.Equal(t, config.StorageBackendBadgerDb, cfg.Storage.Backend.Backend)
		require.Equal(
			t,
			&config.StorageBackendBadgerDbConfig{
				AuthDir:   "./data/auth",
				LogDir:    "/env/logs",
				ConfigDir: "./data/config",
			},
			cfg.Storage.Backend.BadgerDB,
		)
		require.Nil(t, cfg.Storage.Backend.FoundationDB)
	})

	t.Run("env selecting the backend of the file keeps its settings", func(t *testing.T) {
		cfg := loadConfigWithEnv(
			t,
			block("storage", block(`backend "foundationdb"`, `cluster_file = "/file/fdb.cluster"`)),
			map[string]string{"FLOWG_STORAGE_BACKEND": "foundationdb"},
		)

		require.Equal(
			t,
			&config.StorageBackendFoundationDbConfig{
				ClusterFile: "/file/fdb.cluster",
				KeySpace:    "flowg",
			},
			cfg.Storage.Backend.FoundationDB,
		)
		require.Nil(t, cfg.Storage.Backend.BadgerDB)
	})
}
