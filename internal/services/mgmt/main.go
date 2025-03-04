package mgmt

import (
	"log/slog"

	"crypto/tls"
	"net/url"

	"link-society.com/flowg/internal/cluster"
	"link-society.com/flowg/internal/utils/proctree"
)

type ServerOptions struct {
	BindAddress string
	TlsConfig   *tls.Config

	ClusterNodeID       string
	ClusterJoinNodeID   string
	ClusterJoinEndpoint *url.URL
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

		LocalEndpointResolver: func() *url.URL {
			localEndpoint := &url.URL{
				Scheme: "http",
				Host:   listenerH.listener.Addr().String(),
			}

			if opts.TlsConfig != nil {
				localEndpoint.Scheme = "https"
			}

			return localEndpoint
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
