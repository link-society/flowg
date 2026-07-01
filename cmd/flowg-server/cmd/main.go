package cmd

import (
	"fmt"

	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/app/featureflags"
	"link-society.com/flowg/internal/app/metrics"
	"link-society.com/flowg/internal/app/server"

	"link-society.com/flowg/cmd/flowg-server/logging"
)

// ExitCode is the process exit code; the root command sets it to 1 when startup
// or configuration fails.
var ExitCode int = 0

// NewRootCommand builds the root flowg-server command. Its pre-run hook sets the
// umask, configures logging, applies the demo-mode flag and registers the
// Prometheus metrics; its run hook turns the CLI options into a server
// configuration and starts the fx application.
func NewRootCommand() *cobra.Command {
	var configPath string

	rootCmd := &cobra.Command{
		Use:   "flowg-server",
		Short: "Low-Code log management solution",
		Run: func(cmd *cobra.Command, args []string) {
			config := DefaultConfig()

			if configPath != "" {
				if err := config.Load(configPath); err != nil {
					fmt.Printf("ERROR: %v\n", err)
					ExitCode = 1
					return
				}
			}

			opts, err := config.AsServerOptions()
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				ExitCode = 1
				return
			}

			syscall.Umask(0077)
			logging.Setup(config.Logging.Verbose, config.Logging.Level)
			featureflags.SetDemoMode(config.DemoMode)
			metrics.Setup()

			fx.New(fx.NopLogger, server.NewServer(opts)).Run()
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
