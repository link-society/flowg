package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/logging"
)

var exitCode int = 0

var (
	defaultBindAddress = getEnv("FLOWG_BIND_ADDRESS", ":5080")
	defaultAuthDir     = getEnv("FLOWG_AUTH_DIR", "./data/auth")
	defaultConfigDir   = getEnv("FLOWG_CONFIG_DIR", "./data/config")
	defaultLogDir      = getEnv("FLOWG_LOG_DIR", "./data/logs")
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "flowg",
		Short: "Low-Code log management solution",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logging.Discard()
		},
	}

	rootCmd.AddCommand(
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
