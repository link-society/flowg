package config

import (
	"os"

	"strconv"
	"strings"
)

// applyEnv overrides the configuration with the FLOWG_* environment variables
// that are set (an empty variable is treated as unset).
func (cfg *RootConfig) applyEnv() {
	cfg.DemoMode = getEnvBool("FLOWG_DEMO_MODE", cfg.DemoMode)

	cfg.Logging.Level = getEnvString("FLOWG_LOGLEVEL", cfg.Logging.Level)
	cfg.Logging.Verbose = getEnvBool("FLOWG_VERBOSE", cfg.Logging.Verbose)

	cfg.Services.Http.applyEnv()
	cfg.Services.Management.applyEnv()
	cfg.Services.Syslog.applyEnv()

	cfg.Storage.Backend.applyEnv()
	cfg.Storage.Seed.Auth.applyEnv()
}

func (c *HttpConfig) applyEnv() {
	c.BindAddress = getEnvString("FLOWG_HTTP_BIND_ADDRESS", c.BindAddress)
	c.MountPath = getEnvString("FLOWG_HTTP_MOUNT_PATH", c.MountPath)

	c.Tls = applyEnvToggle("FLOWG_HTTP_TLS_ENABLED", c.Tls)
	if c.Tls != nil {
		c.Tls.Cert = getEnvString("FLOWG_HTTP_TLS_CERT", c.Tls.Cert)
		c.Tls.Key = getEnvString("FLOWG_HTTP_TLS_KEY", c.Tls.Key)
	}
}

func (c *ManagementConfig) applyEnv() {
	c.BindAddress = getEnvString("FLOWG_MGMT_BIND_ADDRESS", c.BindAddress)

	c.Tls = applyEnvToggle("FLOWG_MGMT_TLS_ENABLED", c.Tls)
	if c.Tls != nil {
		c.Tls.Cert = getEnvString("FLOWG_MGMT_TLS_CERT", c.Tls.Cert)
		c.Tls.Key = getEnvString("FLOWG_MGMT_TLS_KEY", c.Tls.Key)
	}
}

func (c *SyslogConfig) applyEnv() {
	c.Bind = getEnvString("FLOWG_SYSLOG_BIND_ADDRESS", c.Bind)
	c.Protocol = getEnvString("FLOWG_SYSLOG_PROTOCOL", c.Protocol)
	c.InitialAllowedOrigins = getEnvListString("FLOWG_SYSLOG_INITIAL_ALLOWED_ORIGINS", c.InitialAllowedOrigins)

	c.Tls = applyEnvToggle("FLOWG_SYSLOG_TLS_ENABLED", c.Tls)
	if c.Tls != nil {
		c.Tls.Cert = getEnvString("FLOWG_SYSLOG_TLS_CERT", c.Tls.Cert)
		c.Tls.Key = getEnvString("FLOWG_SYSLOG_TLS_KEY", c.Tls.Key)
		c.Tls.AuthEnabled = getEnvBool("FLOWG_SYSLOG_TLS_AUTH", c.Tls.AuthEnabled)
	}
}

// applyEnv selects the backend named by FLOWG_STORAGE_BACKEND, if set. When it
// differs from the current backend, the current backend settings are dropped in
// favor of the defaults of the selected one. The settings of the selected
// backend are then overridden by their own environment variables.
func (c *StorageBackendConfig) applyEnv() {
	backend := getEnvString("FLOWG_STORAGE_BACKEND", c.Backend)
	if backend != c.Backend {
		c.selectBackend(backend)
	}

	if c.BadgerDB != nil {
		c.BadgerDB.AuthDir = getEnvString("FLOWG_BADGER_AUTH_DIR", c.BadgerDB.AuthDir)
		c.BadgerDB.LogDir = getEnvString("FLOWG_BADGER_LOG_DIR", c.BadgerDB.LogDir)
		c.BadgerDB.ConfigDir = getEnvString("FLOWG_BADGER_CONFIG_DIR", c.BadgerDB.ConfigDir)
	}

	if c.FoundationDB != nil {
		c.FoundationDB.ClusterFile = getEnvString("FLOWG_FOUNDATIONDB_CLUSTER_FILE", c.FoundationDB.ClusterFile)
		c.FoundationDB.ConnectionString = getEnvString("FLOWG_FOUNDATIONDB_CONNECTION_STRING", c.FoundationDB.ConnectionString)
		c.FoundationDB.KeySpace = getEnvString("FLOWG_FOUNDATIONDB_KEY_SPACE", c.FoundationDB.KeySpace)
	}
}

func (c *StorageSeedAuthConfig) applyEnv() {
	c.InitialUser = getEnvString("FLOWG_AUTH_INITIAL_USER", c.InitialUser)
	c.InitialPassword = getEnvString("FLOWG_AUTH_INITIAL_PASSWORD", c.InitialPassword)
	c.ResetUser = getEnvString("FLOWG_AUTH_RESET_USER", c.ResetUser)
	c.ResetPassword = getEnvString("FLOWG_AUTH_RESET_PASSWORD", c.ResetPassword)
}

// applyEnvToggle enables or disables an optional block according to the boolean
// environment variable key: the block is allocated with its defaults when
// enabled and missing, removed when disabled, and left untouched when the
// variable is unset.
func applyEnvToggle[T any](key string, block *T) *T {
	if !getEnvBool(key, block != nil) {
		return nil
	}

	if block == nil {
		block = newWithDefaults[T]()
	}

	return block
}

// getEnvListString reads a comma-separated environment variable into a slice,
// trimming whitespace around each item, or returns defaultValue when unset.
func getEnvListString(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	items := strings.Split(value, ",")
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}

	return items
}

// getEnvString returns the value of the environment variable key, or
// defaultValue when it is unset or empty.
func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvBool parses the environment variable key as a boolean, falling back to
// defaultValue when it is unset and to false when it cannot be parsed.
func getEnvBool(key string, defaultValue bool) bool {
	stringVal := os.Getenv(key)
	if stringVal == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(stringVal)
	if err == nil {
		return value
	}
	return false
}
