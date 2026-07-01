package cmd

import (
	"os"

	"strconv"
	"strings"
)

var (
	defaultDemoMode = getEnvBool("FLOWG_DEMO_MODE", false)

	defaultVerbose  = getEnvBool("FLOWG_VERBOSE", false)
	defaultLogLevel = getEnvString("FLOWG_LOGLEVEL", "info")

	defaultHttpBindAddress = getEnvString("FLOWG_HTTP_BIND_ADDRESS", ":5080")
	defaultHttpMountPath   = getEnvString("FLOWG_HTTP_MOUNT_PATH", "/")
	defaultHttpTlsEnabled  = getEnvBool("FLOWG_HTTP_TLS_ENABLED", false)
	defaultHttpTlsCert     = getEnvString("FLOWG_HTTP_TLS_CERT", "")
	defaultHttpTlsCertKey  = getEnvString("FLOWG_HTTP_TLS_KEY", "")

	defaultMgmtBindAddress = getEnvString("FLOWG_MGMT_BIND_ADDRESS", ":9113")
	defaultMgmtTlsEnabled  = getEnvBool("FLOWG_MGMT_TLS_ENABLED", false)
	defaultMgmtTlsCert     = getEnvString("FLOWG_MGMT_TLS_CERT", "")
	defaultMgmtTlsCertKey  = getEnvString("FLOWG_MGMT_TLS_KEY", "")

	defaultSyslogProtocol = getEnvString("FLOWG_SYSLOG_PROTOCOL", "udp")
	defaultSyslogBindAddr = getEnvString("FLOWG_SYSLOG_BIND_ADDRESS", ":5514")

	defaultSyslogTlsEnabled            = getEnvBool("FLOWG_SYSLOG_TLS_ENABLED", false)
	defaultSyslogTlsCert               = getEnvString("FLOWG_SYSLOG_TLS_CERT", "")
	defaultSyslogTlsCertKey            = getEnvString("FLOWG_SYSLOG_TLS_KEY", "")
	defaultSyslogTlsAuthEnabled        = getEnvBool("FLOWG_SYSLOG_TLS_AUTH", false)
	defaultSyslogInitialAllowedOrigins = getEnvListString("FLOWG_SYSLOG_INITIAL_ALLOWED_ORIGINS", []string{})

	defaultStorageBackend  = getEnvString("FLOWG_STORAGE_BACKEND", "badgerdb")
	defaultBadgerAuthDir   = getEnvString("FLOWG_BADGER_AUTH_DIR", "./data/auth")
	defaultBadgerConfigDir = getEnvString("FLOWG_BADGER_CONFIG_DIR", "./data/config")
	defaultBadgerLogDir    = getEnvString("FLOWG_BADGER_LOG_DIR", "./data/logs")

	defaultAuthInitialUser     = getEnvString("FLOWG_AUTH_INITIAL_USER", "root")
	defaultAuthInitialPassword = getEnvString("FLOWG_AUTH_INITIAL_PASSWORD", "root")

	defaultAuthResetUser     = getEnvString("FLOWG_AUTH_RESET_USER", "")
	defaultAuthResetPassword = getEnvString("FLOWG_AUTH_RESET_PASSWORD", "")
)

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
