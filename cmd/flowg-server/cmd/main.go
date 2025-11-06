package cmd

import (
	"fmt"

	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/app/metrics"
	"link-society.com/flowg/internal/app/server"
)

var ExitCode int = 0

func NewRootCommand() *cobra.Command {
	opts := &options{}

	rootCmd := &cobra.Command{
		Use:   "flowg-server",
		Short: "Low-Code log management solution",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			syscall.Umask(0077)
			logging.Setup(opts.verbose, opts.loglevel)
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
