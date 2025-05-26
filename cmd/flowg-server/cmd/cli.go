package cmd

import "github.com/spf13/cobra"

type options struct {
	httpBindAddress string
	httpTlsEnabled  bool
	httpTlsCert     string
	httpTlsCertKey  string

	mgmtBindAddress string
	mgmtTlsEnabled  bool
	mgmtTlsCert     string
	mgmtTlsCertKey  string

	clusterNodeID       string
	clusterJoinNodeID   string
	clusterJoinEndpoint string
	clusterCookie       string
	clusterStateDir     string

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

	serviceName string
	consulUrl   string

	authInitialUser     string
	authInitialPassword string
}

func (opts *options) defineCliOptions(cmd *cobra.Command) {
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
		&opts.clusterNodeID,
		"cluster-node-id",
		defaultClusterNodeID,
		"Unique identifier for this node in the cluster (leave empty to generate one)",
	)

	cmd.Flags().StringVar(
		&opts.clusterJoinNodeID,
		"cluster-join-node-id",
		defaultClusterJoinNodeID,
		"Unique identifier of the node to join in the cluster",
	)

	cmd.Flags().StringVar(
		&opts.clusterJoinEndpoint,
		"cluster-join-endpoint",
		defaultClusterJoinEndpoint,
		"Management endpoint of the node to join the cluster",
	)

	cmd.Flags().StringVar(
		&opts.clusterCookie,
		"cluster-cookie",
		defaultClusterCookie,
		"Cookie to use for cluster communication (leave empty to disable)",
	)

	cmd.Flags().StringVar(
		&opts.clusterStateDir,
		"cluster-state-dir",
		defaultClusterStateDir,
		"Path to the cluster state directory (for replication)",
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

	cmd.Flags().StringVar(
		&opts.serviceName,
		"service-name",
		defaultServiceName,
		"Name of the service",
	)

	cmd.Flags().StringVar(
		&opts.consulUrl,
		"consul-url",
		defaultConsulUrl,
		"URL to local consul instance",
	)

	cmd.Flags().StringVar(
		&opts.authInitialUser,
		"auth-initial-user",
		defaultAuthInitialUser,
		"Username for the initial admin user (defaults to FLOWG_AUTH_INITIAL_USER or 'root')",
	)

	cmd.Flags().StringVar(
		&opts.authInitialPassword,
		"auth-initial-password",
		defaultAuthInitialPassword,
		"Password for the initial admin user (defaults to FLOWG_AUTH_INITIAL_PASSWORD or 'root')",
	)
}
