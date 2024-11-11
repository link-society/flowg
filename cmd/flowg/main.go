package main

import (
	"fmt"
	"strings"

	"os"
	"syscall"
	"strconv"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/app/logging"
)

var exitCode int = 0

var (
	defaultVerbose = getEnvBool("FLOWG_VERBOSE", false)
	defaultHttpBindAddress = getEnvString("FLOWG_HTTP_BIND_ADDRESS", ":5080")

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

	defaultHttpTlsEnabled = getEnvBool("FLOWG_HTTP_TLS_ENABLED", false)
	defaultHttpTlsCert = getEnvString("FLOWG_HTTP_TLS_CERT", "")
	defaultHttpTlsCertKey = getEnvString("FLOWG_HTTP_TLS_KEY", "")

	defaultSyslogTlsEnabled     = getEnvBool("FLOWG_SYSLOG_TLS_ENABLED", false)
	defaultSyslogTlsCert		= getEnvString("FLOWG_SYSLOG_TLS_CERT", "")
	defaultSyslogTlsCertKey		= getEnvString("FLOWG_SYSLOG_TLS_KEY", "")
	defaultSyslogTlsAuthEnabled	= getEnvBool("FLOWG_SYSLOG_TLS_AUTH", false)

	defaultAuthDir   = getEnvString("FLOWG_AUTH_DIR", "./data/auth")
	defaultConfigDir = getEnvString("FLOWG_CONFIG_DIR", "./data/config")
	defaultLogDir    = getEnvString("FLOWG_LOG_DIR", "./data/logs")
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "flowg",
		Short: "Low-Code log management solution",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logging.Discard()
			syscall.Umask(0077)
		},
	}

	rootCmd.AddCommand(
		NewVersionCommand(),
		NewAdminCommand(),
		NewServeCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
	}

	os.Exit(exitCode)
}

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
