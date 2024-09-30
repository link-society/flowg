package main

import (
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"crypto/tls"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/app/metrics"
	"link-society.com/flowg/internal/app/server"
)

type serveCommandOpts struct {
	httpBindAddress string
	httpTlsEnabled  bool
	httpTlsCert     string
	httpTlsCertKey  string

	syslogBindAddr string

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
			var httpTlsConfig *tls.Config

			if opts.httpTlsEnabled {
				cert, err := tls.LoadX509KeyPair(opts.httpTlsCert, opts.httpTlsCertKey)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to load TLS certificate: %v", err)
					exitCode = 1
					return
				}

				httpTlsConfig = &tls.Config{
					Certificates: []tls.Certificate{cert},
				}
			}

			srv := server.NewServer(server.Options{
				HttpBindAddress: opts.httpBindAddress,
				HttpTlsConfig:   httpTlsConfig,

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

	cmd.Flags().BoolVar(
		&opts.httpTlsEnabled,
		"http-tls",
		false,
		"Enable TLS for the HTTP server",
	)

	cmd.Flags().StringVar(
		&opts.httpTlsCert,
		"http-tls-cert",
		"",
		"Path to the certificate file for the HTTPS server",
	)

	cmd.Flags().StringVar(
		&opts.httpTlsCertKey,
		"http-tls-key",
		"",
		"Path to the certificate key file for the HTTPS server",
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
