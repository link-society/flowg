package cmd

import (
	"os"

	"strconv"
	"strings"
)

var (
	defaultVerbose = getEnvBool("FLOWG_VERBOSE", false)

	defaultHttpBindAddress = getEnvString("FLOWG_HTTP_BIND_ADDRESS", ":5080")
	defaultHttpTlsEnabled  = getEnvBool("FLOWG_HTTP_TLS_ENABLED", false)
	defaultHttpTlsCert     = getEnvString("FLOWG_HTTP_TLS_CERT", "")
	defaultHttpTlsCertKey  = getEnvString("FLOWG_HTTP_TLS_KEY", "")

	defaultMgmtBindAddress = getEnvString("FLOWG_MGMT_BIND_ADDRESS", ":9113")
	defaultMgmtTlsEnabled  = getEnvBool("FLOWG_MGMT_TLS_ENABLED", false)
	defaultMgmtTlsCert     = getEnvString("FLOWG_MGMT_TLS_CERT", "")
	defaultMgmtTlsCertKey  = getEnvString("FLOWG_MGMT_TLS_KEY", "")

	defaultClusterNodeID       = getEnvString("FLOWG_CLUSTER_NODE_ID", "")
	defaultClusterJoinNodeID   = getEnvString("FLOWG_CLUSTER_JOIN_NODE_ID", "")
	defaultClusterJoinEndpoint = getEnvString("FLOWG_CLUSTER_JOIN_ENDPOINT", "")
	defaultClusterCookie       = getEnvString("FLOWG_CLUSTER_COOKIE", "")

	defaultSyslogProtocol     = getEnvString("FLOWG_SYSLOG_PROTOCOL", "udp")
	defaultSyslogBindAddr     = getEnvString("FLOWG_SYSLOG_BIND_ADDRESS", ":5514")
	defaultSyslogAllowOrigins = (func() []string {
		origins := getEnvString("FLOWG_SYSLOG_ALLOW_ORIGINS", "")
		if origins == "" {
			return nil
		} else {
			return strings.Split(origins, ",")
		}
	})()

	defaultSyslogTlsEnabled     = getEnvBool("FLOWG_SYSLOG_TLS_ENABLED", false)
	defaultSyslogTlsCert        = getEnvString("FLOWG_SYSLOG_TLS_CERT", "")
	defaultSyslogTlsCertKey     = getEnvString("FLOWG_SYSLOG_TLS_KEY", "")
	defaultSyslogTlsAuthEnabled = getEnvBool("FLOWG_SYSLOG_TLS_AUTH", false)

	defaultAuthDir         = getEnvString("FLOWG_AUTH_DIR", "./data/auth")
	defaultConfigDir       = getEnvString("FLOWG_CONFIG_DIR", "./data/config")
	defaultLogDir          = getEnvString("FLOWG_LOG_DIR", "./data/logs")
	defaultClusterStateDir = getEnvString("FLOWG_CLUSTER_STATE_DIR", "./data/state")

	defaultServiceName = getEnvString("FLOWG_SERVICE_NAME", "FlowG")
	defaultConsulUrl   = getEnvString("FLOWG_CONSUL_URL", "")

	defaultAuthInitialUser     = getEnvString("FLOWG_AUTH_INITIAL_USER", "root")
	defaultAuthInitialPassword = getEnvString("FLOWG_AUTH_INITIAL_PASSWORD", "root")

	defaultAuthResetUser     = getEnvString("FLOWG_AUTH_RESET_USER", "")
	defaultAuthResetPassword = getEnvString("FLOWG_AUTH_RESET_PASSWORD", "")
)

func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

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
