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
	opts := &options{}

	rootCmd := &cobra.Command{
		Use:   "flowg-server",
		Short: "Low-Code log management solution",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			syscall.Umask(0077)
			logging.Setup(opts.verbose, opts.loglevel)
			featureflags.SetDemoMode(opts.demoMode)
			metrics.Setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			opts, err := newServerConfig(opts)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				ExitCode = 1
				return
			}

			fx.New(fx.NopLogger, server.NewServer(opts)).Run()
		},
	}

	opts.defineCliOptions(rootCmd)

	return rootCmd
}
