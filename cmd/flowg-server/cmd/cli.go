package cmd

import (
	"github.com/spf13/cobra"
)

type options struct {
	demoMode bool

	httpBindAddress string
	httpMountPath   string
	httpTlsEnabled  bool
	httpTlsCert     string
	httpTlsCertKey  string

	mgmtBindAddress string
	mgmtTlsEnabled  bool
	mgmtTlsCert     string
	mgmtTlsCertKey  string

	syslogProtocol              string
	syslogBindAddr              string
	syslogTlsEnabled            bool
	syslogTlsCert               string
	syslogTlsCertKey            string
	syslogTlsAuthEnabled        bool
	syslogInitialAllowedOrigins []string

	authDir   string
	logDir    string
	configDir string

	verbose  bool
	loglevel string

	authInitialUser     string
	authInitialPassword string

	authResetUser     string
	authResetPassword string
}

func (opts *options) defineCliOptions(cmd *cobra.Command) {
	cmd.Flags().BoolVar(
		&opts.demoMode,
		"demo-mode",
		defaultDemoMode,
		"Enable demo mode (add demo account and limit storage usage)",
	)

	cmd.Flags().StringVar(
		&opts.httpBindAddress,
		"http-bind",
		defaultHttpBindAddress,
		"Address to bind the HTTP server to",
	)

	cmd.Flags().StringVar(
		&opts.httpMountPath,
		"http-mount-path",
		defaultHttpMountPath,
		"Path to mount the HTTP server",
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
		&opts.syslogInitialAllowedOrigins,
		"syslog-initial-allowed-origin",
		defaultSyslogInitialAllowedOrigins,
		"Initial allowed origins for the Syslog server (can be set multiple times)",
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

	cmd.Flags().StringVar(
		&opts.loglevel,
		"loglevel",
		defaultLogLevel,
		"Set the logging level (one of \"debug\", \"info\", \"warn\", \"error\"), ignored if 'verbose' is set",
	)

	cmd.Flags().StringVar(
		&opts.authInitialUser,
		"auth-initial-user",
		defaultAuthInitialUser,
		"Username for the initial admin user",
	)

	cmd.Flags().StringVar(
		&opts.authInitialPassword,
		"auth-initial-password",
		defaultAuthInitialPassword,
		"Password for the initial admin user",
	)

	cmd.Flags().StringVar(
		&opts.authResetUser,
		"auth-reset-user",
		defaultAuthResetUser,
		"If set, this is the username for the user to reset the password for",
	)

	cmd.Flags().StringVar(
		&opts.authResetPassword,
		"auth-reset-password",
		defaultAuthResetPassword,
		"If set, this is the new password for the user to reset",
	)
}
