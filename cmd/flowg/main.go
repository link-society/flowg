package main

import (
	"fmt"

	"os"
	"syscall"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/app/logging"
)

var exitCode int = 0

var (
	defaultHttpBindAddress = getEnv("FLOWG_HTTP_BIND_ADDRESS", ":5080")

	defaultSyslogProtocol = getEnv("FLOWG_SYSLOG_PROTOCOL", "udp")
	defaultSyslogBindAddr = getEnv("FLOWG_SYSLOG_BIND_ADDRESS", ":5514")

	defaultAuthDir   = getEnv("FLOWG_AUTH_DIR", "./data/auth")
	defaultConfigDir = getEnv("FLOWG_CONFIG_DIR", "./data/config")
	defaultLogDir    = getEnv("FLOWG_LOG_DIR", "./data/logs")
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

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
