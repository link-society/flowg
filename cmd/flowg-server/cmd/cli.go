package cmd

import (
	"github.com/spf13/cobra"
)

type options struct {
	httpBindAddress string
	httpTlsEnabled  bool
	httpTlsCert     string
	httpTlsCertKey  string

	mgmtBindAddress string
	mgmtTlsEnabled  bool
	mgmtTlsCert     string
	mgmtTlsCertKey  string

	clusterNodeID   string
	clusterCookie   string
	clusterStateDir string

	clusterFormationStrategy string

	clusterFormationManualJoinNodeID   string
	clusterFormationManualJoinEndpoint string

	clusterFormationConsulServiceName string
	clusterFormationConsulUrl         string

	clusterFormationKubernetesServiceNamespace string
	clusterFormationKubernetesServiceName      string
	clusterFormationKubernetesServicePortName  string

	clusterFormationDnsServiceName   string
	clusterFormationDnsServerAddress string
	clusterFormationDnsDomainName    string

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

	verbose  bool
	loglevel string

	authInitialUser     string
	authInitialPassword string

	authResetUser     string
	authResetPassword string
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
		&opts.clusterFormationStrategy,
		"cluster-formation-strategy",
		defaultClusterFormationStrategy,
		"Strategy to use for cluster formation",
	)

	cmd.Flags().StringVar(
		&opts.clusterFormationManualJoinNodeID,
		"cluster-formation-manual-join-node-id",
		defaultClusterFormationManualJoinNodeID,
		"Unique identifier of the node to join in the cluster, ignored if not using 'manual' strategy",
	)

	cmd.Flags().StringVar(
		&opts.clusterFormationManualJoinEndpoint,
		"cluster-formation-manual-join-endpoint",
		defaultClusterFormationManualJoinEndpoint,
		"Management endpoint of the node to join the cluster, ignored if not using 'manual' strategy",
	)

	cmd.Flags().StringVar(
		&opts.clusterFormationConsulServiceName,
		"cluster-formation-consul-service-name",
		defaultClusterFormationConsulServiceName,
		"Name of the Consul service, ignored if not using 'consul' strategy",
	)

	cmd.Flags().StringVar(
		&opts.clusterFormationConsulUrl,
		"cluster-formation-consul-url",
		defaultClusterFormationConsulUrl,
		"URL to local consul instance, ignored if not using 'consul' strategy",
	)

	cmd.Flags().StringVar(
		&opts.clusterFormationKubernetesServiceNamespace,
		"cluster-formation-k8s-service-namespace",
		defaultClusterFormationKubernetesServiceNamespace,
		"Namespace of the Kubernetes service, ignored if not using 'k8s' strategy",
	)

	cmd.Flags().StringVar(
		&opts.clusterFormationKubernetesServiceName,
		"cluster-formation-k8s-service-name",
		defaultClusterFormationKubernetesServiceName,
		"Name of the Kubernetes service, ignored if not using 'k8s' strategy",
	)

	cmd.Flags().StringVar(
		&opts.clusterFormationKubernetesServicePortName,
		"cluster-formation-k8s-service-port-name",
		defaultClusterFormationKubernetesServicePortName,
		"Name of the port in the Kubernetes service to use for cluster formation, ignored if not using 'k8s' strategy",
	)

	cmd.Flags().StringVar(
		&opts.clusterFormationDnsServiceName,
		"cluster-formation-dns-service-name",
		defaultClusterFormationDnsServiceName,
		"Name of the dns service to use for cluster formation, ignored if not using 'dns' strategy",
	)

	cmd.Flags().StringVar(
		&opts.clusterFormationDnsDomainName,
		"cluster-formation-dns-domain-name",
		defaultClusterFormationDnsDomainName,
		"Domain name of the your domain to use for cluster formation, ignored if not using 'dns' strategy",
	)
	//DnsServerAddress: opts.clusterFormationDnsServerAddress,

	cmd.Flags().StringVar(
		&opts.clusterFormationDnsServerAddress,
		"cluster-formation-dns-server-address",
		defaultClusterFormationDnsServerAddress,
		"DNS Server address to use for cluster formation, ignored if not using 'dns' strategy",
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
