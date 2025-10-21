package cmd

import (
	"fmt"

	"strings"

	"crypto/tls"
	"net"
	"net/url"

	"link-society.com/flowg/internal/app/server"
	"link-society.com/flowg/internal/cluster"
	"link-society.com/flowg/internal/utils/rnd"
)

func newServerConfig(opts *options) (server.Options, error) {
	var (
		httpTlsConfig   *tls.Config
		mgmtTlsConfig   *tls.Config
		syslogTlsConfig *tls.Config
	)

	if opts.httpTlsEnabled {
		cert, err := tls.LoadX509KeyPair(opts.httpTlsCert, opts.httpTlsCertKey)
		if err != nil {
			return server.Options{}, fmt.Errorf("failed to load TLS certificate: %w", err)
		}

		httpTlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	if opts.mgmtTlsEnabled {
		cert, err := tls.LoadX509KeyPair(opts.mgmtTlsCert, opts.mgmtTlsCertKey)
		if err != nil {
			return server.Options{}, fmt.Errorf("failed to load Management TLS certificate: %w", err)
		}

		mgmtTlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	if opts.clusterNodeID == "" {
		opts.clusterNodeID = rnd.RandomName()
	}

	var clusterFormationStrategy cluster.ClusterFormationStrategy

	switch opts.clusterFormationStrategy {
	case "manual":
		manualStrategy := &cluster.ManualClusterFormationStrategy{}

		if opts.clusterFormationManualJoinNodeID != "" {
			if opts.clusterFormationManualJoinEndpoint == "" {
				return server.Options{}, fmt.Errorf("cluster join endpoint is required when joining a cluster")
			}

			endpoint, err := url.Parse(opts.clusterFormationManualJoinEndpoint)
			if err != nil {
				return server.Options{}, fmt.Errorf("invalid cluster join endpoint: %w", err)
			}

			manualStrategy.JoinNodeID = opts.clusterFormationManualJoinNodeID
			manualStrategy.JoinNodeEndpoint = endpoint
		}

		clusterFormationStrategy = manualStrategy

	case "consul":
		if opts.clusterFormationConsulServiceName == "" {
			return server.Options{}, fmt.Errorf("service name is required for 'consul' cluster formation")
		}
		if opts.clusterFormationConsulUrl == "" {
			return server.Options{}, fmt.Errorf("consul URL is required for 'consul' cluster formation")
		}

		clusterFormationStrategy = &cluster.ConsulClusterFormationStrategy{
			NodeID:         opts.clusterNodeID,
			ServiceName:    opts.clusterFormationConsulServiceName,
			ServiceAddress: opts.mgmtBindAddress,
			ServiceTls:     opts.mgmtTlsEnabled,
			ConsulUrl:      opts.clusterFormationConsulUrl,
		}

	case "k8s":
		if opts.clusterFormationKubernetesServiceNamespace == "" {
			return server.Options{}, fmt.Errorf("kubernetes service namespace is required for 'k8s' cluster formation")
		}
		if opts.clusterFormationKubernetesServiceName == "" {
			return server.Options{}, fmt.Errorf("kubernetes service name is required for 'k8s' cluster formation")
		}
		if opts.clusterFormationKubernetesServicePortName == "" {
			return server.Options{}, fmt.Errorf("kubernetes service port name is required for 'k8s' cluster formation")
		}

		clusterFormationStrategy = &cluster.KubernetesClusterFormationStrategy{
			NodeID:           opts.clusterNodeID,
			ServiceNamespace: opts.clusterFormationKubernetesServiceNamespace,
			ServiceName:      opts.clusterFormationKubernetesServiceName,
			ServicePortName:  opts.clusterFormationKubernetesServicePortName,
		}

	case "dns":
		if opts.clusterFormationDnsServiceName == "" {
			return server.Options{}, fmt.Errorf("dns service name is required for 'dns' cluster formation")
		}
		if opts.clusterFormationDnsDomainName == "" {
			return server.Options{}, fmt.Errorf("dns domain name is required for 'dns' cluster formation")
		}
		if opts.clusterFormationDnsServerAddress == "" {
			return server.Options{}, fmt.Errorf("dns server address is required for 'dns' cluster formation")
		}

		clusterFormationStrategy = &cluster.DnsClusterFormationStrategy{
			NodeID:           opts.clusterNodeID,
			ServiceName:      opts.clusterFormationDnsServiceName,
			DnsDomainName:    opts.clusterFormationDnsDomainName,
			DnsServerAddress: opts.clusterFormationDnsServerAddress,
			ServiceAddress:   opts.mgmtBindAddress,
		}

	default:
		return server.Options{}, fmt.Errorf("invalid cluster formation strategy: %s", opts.clusterFormationStrategy)
	}

	if opts.syslogProtocol != "tcp" && opts.syslogProtocol != "udp" {
		return server.Options{}, fmt.Errorf("invalid syslog protocol: %s", opts.syslogProtocol)
	}

	if opts.syslogTlsEnabled && opts.syslogProtocol == "udp" {
		return server.Options{}, fmt.Errorf("TLS is not supported for Syslog UDP protocol")
	}

	if opts.syslogTlsEnabled {
		cert, err := tls.LoadX509KeyPair(opts.syslogTlsCert, opts.syslogTlsCertKey)
		if err != nil {
			return server.Options{}, fmt.Errorf("failed to load Syslog TLS certificate: %w", err)
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
					return server.Options{}, fmt.Errorf("invalid syslog allow origin: %s", origin)
				}
			} else {
				if net.ParseIP(origin) == nil {
					return server.Options{}, fmt.Errorf("invalid syslog allow origin: %s", origin)
				}
			}
		}
	}

	config := server.Options{
		HttpBindAddress: opts.httpBindAddress,
		HttpTlsConfig:   httpTlsConfig,

		MgmtBindAddress: opts.mgmtBindAddress,
		MgmtTlsConfig:   mgmtTlsConfig,

		ClusterNodeID:            opts.clusterNodeID,
		ClusterCookie:            opts.clusterCookie,
		ClusterStateDir:          opts.clusterStateDir,
		ClusterFormationStrategy: clusterFormationStrategy,

		SyslogTcpMode:      opts.syslogProtocol == "tcp",
		SyslogBindAddress:  opts.syslogBindAddr,
		SyslogTlsConfig:    syslogTlsConfig,
		SyslogAllowOrigins: opts.syslogAllowOrigins,

		ConfigStorageDir: opts.configDir,
		AuthStorageDir:   opts.authDir,
		LogStorageDir:    opts.logDir,

		AuthInitialUser:     opts.authInitialUser,
		AuthInitialPassword: opts.authInitialPassword,

		AuthResetUser:     opts.authResetUser,
		AuthResetPassword: opts.authResetPassword,
	}

	return config, nil
}
