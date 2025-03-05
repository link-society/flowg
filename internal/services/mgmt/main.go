package mgmt

import (
	"fmt"
	"log/slog"
	"net"

	"crypto/tls"
	"net/url"

	"github.com/hashicorp/go-sockaddr"

	"link-society.com/flowg/internal/cluster"

	"link-society.com/flowg/internal/utils/proctree"
)

type ServerOptions struct {
	BindAddress string
	TlsConfig   *tls.Config

	ClusterNodeID       string
	ClusterJoinNodeID   string
	ClusterJoinEndpoint *url.URL
	ClusterCookie       string
}

func NewServer(opts *ServerOptions) proctree.Process {
	logger := slog.Default().With(
		slog.String("channel", "mgmt"),
		slog.Group("mgmt",
			slog.String("bind", opts.BindAddress),
		),
	)

	listenerH := &listenerHandler{
		logger:      logger,
		bindAddress: opts.BindAddress,
	}

	clusterManager := cluster.NewManager(&cluster.ManagerOptions{
		NodeID:           opts.ClusterNodeID,
		JoinNodeID:       opts.ClusterJoinNodeID,
		JoinNodeEndpoint: opts.ClusterJoinEndpoint,
		Cookie:           opts.ClusterCookie,

		LocalEndpointResolver: func() (*url.URL, error) {
			host, port, err := net.SplitHostPort(listenerH.listener.Addr().String())
			if err != nil {
				return nil, fmt.Errorf("failed to parse listener address: %w", err)
			}

			if host == "0.0.0.0" || host == "::" {
				ip, err := sockaddr.GetPrivateIP()
				if err != nil {
					return nil, fmt.Errorf("failed to get private IP: %w", err)
				}
				if ip == "" {
					return nil, fmt.Errorf("no private IP found")
				}

				host = ip
			}

			localEndpoint := &url.URL{
				Scheme: "http",
				Host:   net.JoinHostPort(host, port),
			}

			if opts.TlsConfig != nil {
				localEndpoint.Scheme = "https"
			}

			return localEndpoint, nil
		},
	})

	serverH := &serverHandler{
		logger: logger,

		bindAddress: opts.BindAddress,
		tlsConfig:   opts.TlsConfig,

		listenerH:      listenerH,
		clusterManager: clusterManager,
	}

	return proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewProcess(listenerH),
		clusterManager,
		proctree.NewProcess(serverH),
	)
}
