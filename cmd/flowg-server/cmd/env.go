package cmd

import (
	"os"

	"strconv"
	"strings"
)

var (
	defaultVerbose  = getEnvBool("FLOWG_VERBOSE", false)
	defaultLogLevel = getEnvString("FLOWG_LOGLEVEL", "info")

	defaultHttpBindAddress = getEnvString("FLOWG_HTTP_BIND_ADDRESS", ":5080")
	defaultHttpTlsEnabled  = getEnvBool("FLOWG_HTTP_TLS_ENABLED", false)
	defaultHttpTlsCert     = getEnvString("FLOWG_HTTP_TLS_CERT", "")
	defaultHttpTlsCertKey  = getEnvString("FLOWG_HTTP_TLS_KEY", "")

	defaultMgmtBindAddress = getEnvString("FLOWG_MGMT_BIND_ADDRESS", ":9113")
	defaultMgmtTlsEnabled  = getEnvBool("FLOWG_MGMT_TLS_ENABLED", false)
	defaultMgmtTlsCert     = getEnvString("FLOWG_MGMT_TLS_CERT", "")
	defaultMgmtTlsCertKey  = getEnvString("FLOWG_MGMT_TLS_KEY", "")

	defaultClusterNodeID = getEnvString("FLOWG_CLUSTER_NODE_ID", "")
	defaultClusterCookie = getEnvString("FLOWG_CLUSTER_COOKIE", "")

	defaultClusterFormationStrategy = getEnvString("FLOWG_CLUSTER_FORMATION_STRATEGY", "manual")

	defaultClusterFormationManualJoinNodeID   = getEnvString("FLOWG_CLUSTER_FORMATION_MANUAL_JOIN_NODE_ID", "")
	defaultClusterFormationManualJoinEndpoint = getEnvString("FLOWG_CLUSTER_FORMATION_MANUAL_JOIN_ENDPOINT", "")

	defaultClusterFormationConsulServiceName = getEnvString("FLOWG_CLUSTER_FORMATION_CONSUL_SERVICE_NAME", "FlowG")
	defaultClusterFormationConsulUrl         = getEnvString("FLOWG_CLUSTER_FORMATION_CONSUL_URL", "")

	defaultClusterFormationKubernetesServiceNamespace = getEnvString("FLOWG_CLUSTER_FORMATION_K8S_SERVICE_NAMESPACE", "default")
	defaultClusterFormationKubernetesServiceName      = getEnvString("FLOWG_CLUSTER_FORMATION_K8S_SERVICE_NAME", "flowg")
	defaultClusterFormationKubernetesServicePortName  = getEnvString("FLOWG_CLUSTER_FORMATION_K8S_SERVICE_PORT_NAME", "mgmt")

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
