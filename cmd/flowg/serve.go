package main

import (
	"context"
	"time"

	"fmt"
	"strings"

	"os"
	"os/signal"
	"syscall"

	"crypto/tls"
	"net"

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

	mgmtBindAddress string
	mgmtTlsEnabled  bool
	mgmtTlsCert     string
	mgmtTlsCertKey  string

	syslogProtocol       string
	syslogBindAddr       string
	syslogTlsEnabled     bool
	syslogTlsCert        string
	syslogTlsCertKey     string
	syslogTlsAuthEnabled bool
	syslogAllowOrigins   []string

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
			var (
				httpTlsConfig   *tls.Config
				mgmtTlsConfig   *tls.Config
				syslogTlsConfig *tls.Config
			)

			if opts.httpTlsEnabled {
				cert, err := tls.LoadX509KeyPair(opts.httpTlsCert, opts.httpTlsCertKey)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to load TLS certificate: %v\n", err)
					exitCode = 1
					return
				}

				httpTlsConfig = &tls.Config{
					Certificates: []tls.Certificate{cert},
				}
			}

			if opts.mgmtTlsEnabled {
				cert, err := tls.LoadX509KeyPair(opts.mgmtTlsCert, opts.mgmtTlsCertKey)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to load Management TLS certificate: %v\n", err)
					exitCode = 1
					return
				}

				mgmtTlsConfig = &tls.Config{
					Certificates: []tls.Certificate{cert},
				}
			}

			if opts.syslogProtocol != "tcp" && opts.syslogProtocol != "udp" {
				cmd.Usage()
				fmt.Fprintf(os.Stderr, "\nERROR: Invalid syslog protocol: %s\n", opts.syslogProtocol)
				exitCode = 1
				return
			}

			if opts.syslogTlsEnabled && opts.syslogProtocol == "udp" {
				cmd.Usage()
				fmt.Fprintf(os.Stderr, "\nERROR: TLS is not supported for Syslog UDP protocol\n")
				exitCode = 1
				return
			}

			if opts.syslogTlsEnabled {
				cert, err := tls.LoadX509KeyPair(opts.syslogTlsCert, opts.syslogTlsCertKey)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to load Syslog TLS certificate: %v\n", err)
					exitCode = 1
					return
				}

				clientAuth := tls.VerifyClientCertIfGiven
				if opts.syslogTlsAuthEnabled {
					clientAuth = tls.RequireAndVerifyClientCert
				}

				syslogTlsConfig = &tls.Config{
					Certificates: []tls.Certificate{cert},
					ClientAuth:   clientAuth,
				}
			}

			if opts.syslogAllowOrigins != nil {
				for _, origin := range opts.syslogAllowOrigins {
					if strings.Contains(origin, "/") {
						_, _, err := net.ParseCIDR(origin)
						if err != nil {
							cmd.Usage()
							fmt.Fprintf(os.Stderr, "\nERROR: Invalid syslog allow origin: %s\n", origin)
							exitCode = 1
							return
						}
					} else {
						if net.ParseIP(origin) == nil {
							cmd.Usage()
							fmt.Fprintf(os.Stderr, "\nERROR: Invalid syslog allow origin: %s\n", origin)
							exitCode = 1
							return
						}
					}
				}
			}

			srv := server.NewServer(server.Options{
				HttpBindAddress: opts.httpBindAddress,
				HttpTlsConfig:   httpTlsConfig,

				MgmtBindAddress: opts.mgmtBindAddress,
				MgmtTlsConfig:   mgmtTlsConfig,

				SyslogTCP:          opts.syslogProtocol == "tcp",
				SyslogBindAddress:  opts.syslogBindAddr,
				SyslogTlsConfig:    syslogTlsConfig,
				SyslogAllowOrigins: opts.syslogAllowOrigins,

				ConfigStorageDir: opts.configDir,
				AuthStorageDir:   opts.authDir,
				LogStorageDir:    opts.logDir,
			})

			srv.Start()
			err := srv.WaitReady(context.Background())
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to start server: %v\n", err)
				exitCode = 1
				return
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
					fmt.Fprintf(os.Stderr, "ERROR: Failed to stop server: %v\n", err)
					exitCode = 1
					return
				}

			case err := <-doneC:
				monitorCancel()
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Server stopped unexpectedly: %v\n", err)
					exitCode = 1
					return
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
		defaultHttpTlsEnabled,
		"Enable TLS for the HTTP server",
	)

	cmd.Flags().StringVar(
		&opts.httpTlsCert,
		"http-tls-cert",
		defaultHttpTlsCert,
		"Path to the certificate file for the HTTPS server",
	)

	cmd.Flags().StringVar(
		&opts.httpTlsCertKey,
		"http-tls-key",
		defaultHttpTlsCertKey,
		"Path to the certificate key file for the HTTPS server",
	)

	cmd.Flags().StringVar(
		&opts.mgmtBindAddress,
		"mgmt-bind",
		defaultMgmtBindAddress,
		"Address to bind the Management HTTP server to",
	)

	cmd.Flags().BoolVar(
		&opts.mgmtTlsEnabled,
		"mgmt-tls",
		defaultMgmtTlsEnabled,
		"Enable TLS for the Management HTTP server",
	)

	cmd.Flags().StringVar(
		&opts.mgmtTlsCert,
		"mgmt-tls-cert",
		defaultMgmtTlsCert,
		"Path to the certificate file for the Management HTTPS server",
	)

	cmd.Flags().StringVar(
		&opts.mgmtTlsCertKey,
		"mgmt-tls-key",
		defaultMgmtTlsCertKey,
		"Path to the certificate key file for the Management HTTPS server",
	)

	cmd.Flags().StringVar(
		&opts.syslogProtocol,
		"syslog-proto",
		defaultSyslogProtocol,
		"Protocol to use for the Syslog server (one of \"tcp\" or \"udp\")",
	)

	cmd.Flags().StringVar(
		&opts.syslogBindAddr,
		"syslog-bind",
		defaultSyslogBindAddr,
		"Address to bind the Syslog server to",
	)

	cmd.Flags().BoolVar(
		&opts.syslogTlsEnabled,
		"syslog-tls",
		defaultSyslogTlsEnabled,
		"Enable TLS for the Syslog server (requires protocol to be \"tcp\")",
	)

	cmd.Flags().StringVar(
		&opts.syslogTlsCert,
		"syslog-tls-cert",
		defaultSyslogTlsCert,
		"Path to the certificate file for the Syslog server",
	)

	cmd.Flags().StringVar(
		&opts.syslogTlsCertKey,
		"syslog-tls-key",
		defaultSyslogTlsCertKey,
		"Path to the certificate key file for the Syslog server",
	)

	cmd.Flags().BoolVar(
		&opts.syslogTlsAuthEnabled,
		"syslog-tls-auth",
		defaultSyslogTlsAuthEnabled,
		"Require clients to authenticate against the Syslog server with a client certificate",
	)

	cmd.Flags().StringArrayVar(
		&opts.syslogAllowOrigins,
		"syslog-allow-origin",
		defaultSyslogAllowOrigins,
		"Allowed origin (IP address or CIDR range) for Syslog server (default: all)",
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
		defaultVerbose,
		"Enable verbose logging",
	)

	return cmd
}
