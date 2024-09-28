package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/app/metrics"
	"link-society.com/flowg/internal/app/server"
)

type serveCommandOpts struct {
	httpBindAddress string
	syslogBindAddr  string

	authDir   string
	logDir    string
	configDir string
	verbose   bool
}

func NewServeCommand() *cobra.Command {
	opts := &serveCommandOpts{}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start FlowG standalone server",
		PreRun: func(cmd *cobra.Command, args []string) {
			logging.Setup(opts.verbose)
			metrics.Setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			srv := server.NewServer(server.Options{
				HttpBindAddress:   opts.httpBindAddress,
				SyslogBindAddress: opts.syslogBindAddr,

				ConfigStorageDir: opts.configDir,
				AuthStorageDir:   opts.authDir,
				LogStorageDir:    opts.logDir,
			})

			srv.Start()

			sigC := make(chan os.Signal, 1)
			signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

			select {
			case <-sigC:
				srv.Stop()
				failure := <-srv.DoneC()
				if failure {
					exitCode = 1
				}

			case failure := <-srv.DoneC():
				if failure {
					exitCode = 1
				}
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.httpBindAddress,
		"http-bind",
		defaultHttpBindAddress,
		"Address to bind the HTTP server to",
	)

	cmd.Flags().StringVar(
		&opts.syslogBindAddr,
		"syslog-bind",
		defaultSyslogBindAddr,
		"Address to bind the Syslog server to",
	)

	cmd.Flags().StringVar(
		&opts.authDir,
		"auth-dir",
		defaultAuthDir,
		"Path to the auth database directory",
	)

	cmd.Flags().StringVar(
		&opts.logDir,
		"log-dir",
		defaultLogDir,
		"Path to the log database directory",
	)
	cmd.MarkFlagDirname("log-dir")

	cmd.Flags().StringVar(
		&opts.configDir,
		"config-dir",
		defaultConfigDir,
		"Path to the config directory",
	)
	cmd.MarkFlagDirname("config-dir")

	cmd.Flags().BoolVar(
		&opts.verbose,
		"verbose",
		false,
		"Enable verbose logging",
	)

	return cmd
}
