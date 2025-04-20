package cmd

import (
	"context"
	"time"

	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

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
			logging.Setup(opts.verbose)
			metrics.Setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				config, err := newServerConfig(opts)
				if err != nil {
					return fmt.Errorf("failed to create server configuration: %w", err)
				}

				srv := server.NewServer(config)

				srv.Start()
				if err := srv.WaitReady(context.Background()); err != nil {
					return fmt.Errorf("failed to start server: %w", err)
				}

				monitorCtx, monitorCancel := context.WithCancel(context.Background())
				doneC := make(chan error, 1)
				go func() {
					err := srv.Join(monitorCtx)
					doneC <- err
				}()

				sigC := make(chan os.Signal, 1)
				signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

				select {
				case <-sigC:
					monitorCancel()
					srv.Stop()

					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					err := srv.Join(ctx)
					if err != nil {
						return fmt.Errorf("failed to stop server: %w", err)
					}

				case err := <-doneC:
					monitorCancel()
					if err != nil {
						return fmt.Errorf("server stopped unexpectedly: %w", err)
					}
				}

				return nil
			}()

			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				ExitCode = 1
			}
		},
	}

	opts.defineCliOptions(rootCmd)

	return rootCmd
}
