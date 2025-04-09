package cmd

import (
	"fmt"

	"strings"

	"crypto/tls"
	"net"
	"net/url"

	"link-society.com/flowg/internal/app/server"
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

	var clusterJoinEndpointUrl *url.URL
	if opts.clusterJoinNodeID != "" {
		if opts.clusterJoinEndpoint == "" {
			return server.Options{}, fmt.Errorf("cluster join endpoint is required when joining a cluster")
		}

		var err error
		clusterJoinEndpointUrl, err = url.Parse(opts.clusterJoinEndpoint)
		if err != nil {
			return server.Options{}, fmt.Errorf("invalid cluster join endpoint: %w", err)
		}
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

		ClusterNodeID:       opts.clusterNodeID,
		ClusterJoinNodeID:   opts.clusterJoinNodeID,
		ClusterJoinEndpoint: clusterJoinEndpointUrl,
		ClusterCookie:       opts.clusterCookie,

		SyslogTcpMode:      opts.syslogProtocol == "tcp",
		SyslogBindAddress:  opts.syslogBindAddr,
		SyslogTlsConfig:    syslogTlsConfig,
		SyslogAllowOrigins: opts.syslogAllowOrigins,

		ConfigStorageDir: opts.configDir,
		AuthStorageDir:   opts.authDir,
		LogStorageDir:    opts.logDir,

		ConsulUrl: opts.consulUrl,
	}

	return config, nil
}
