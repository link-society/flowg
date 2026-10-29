package cmd

import (
	"fmt"

	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"go.uber.org/fx"

	"go.link-society.com/flowg/internal/app/config"
	"go.link-society.com/flowg/internal/app/metrics"
	"go.link-society.com/flowg/internal/app/server"

	"go.link-society.com/flowg/cmd/flowg-server/logging"
)

// ExitCode is the process exit code; the root command sets it to 1 when startup
// or configuration fails.
var ExitCode int = 0

// NewRootCommand builds the root flowg-server command. Its run hook loads the
// configuration, sets the umask, configures logging, registers the Prometheus
// metrics, and starts the fx application built from the configuration.
func NewRootCommand() *cobra.Command {
	var configPath string

	rootCmd := &cobra.Command{
		Use:   "flowg-server",
		Short: "Low-Code log management solution",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.DefaultConfig()
			if err := cfg.Load(configPath); err != nil {
				fmt.Printf("ERROR: %v\n", err)
				ExitCode = 1
				return
			}

			syscall.Umask(0077)
			logger := logging.Setup(cfg.Logging.Verbose, cfg.Logging.Level)
			metrics.Setup()

			app := fx.New(fx.NopLogger, server.NewServer(cfg, logger))
			if err := app.Err(); err != nil {
				fmt.Printf("ERROR: %v\n", dig.RootCause(err))
				ExitCode = 1
				return
			}

			app.Run()
		},
	}

	rootCmd.Flags().StringVar(
		&configPath,
		"config",
		"",
		"Path to the configuration file",
	)

	return rootCmd
}
