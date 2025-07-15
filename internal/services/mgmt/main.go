package mgmt

import (
	"fmt"
	"log/slog"
	"time"

	"crypto/tls"
	"net"
	"net/url"

	"github.com/hashicorp/go-sockaddr"

	"link-society.com/flowg/internal/cluster"
	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/utils/kvstore"
	"link-society.com/flowg/internal/utils/proctree"
)

type ServerOptions struct {
	BindAddress string
	TlsConfig   *tls.Config

	ClusterNodeID            string
	ClusterCookie            string
	ClusterStateDir          string
	ClusterFormationStrategy cluster.ClusterFormationStrategy

	AuthStorage   *auth.Storage
	ConfigStorage *config.Storage
	LogStorage    *log.Storage
}

func NewServer(opts *ServerOptions) proctree.Process {
	logger := slog.Default().With(
		slog.String("channel", "mgmt"),
		slog.Group("mgmt",
			slog.String("bind", opts.BindAddress),
		),
	)

	clusterStateStorage := kvstore.NewStorage(kvstore.OptDirectory(opts.ClusterStateDir))

	listenerH := &listenerHandler{
		logger:      logger,
		bindAddress: opts.BindAddress,
	}

	clusterManager := cluster.NewManager(&cluster.ManagerOptions{
		NodeID: opts.ClusterNodeID,
		Cookie: opts.ClusterCookie,

		ClusterFormationStrategy: opts.ClusterFormationStrategy,

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

		AuthStorage:         opts.AuthStorage,
		ConfigStorage:       opts.ConfigStorage,
		LogStorage:          opts.LogStorage,
		ClusterStateStorage: clusterStateStorage,
	})

	serverH := &serverHandler{
		logger: logger,

		bindAddress: opts.BindAddress,
		tlsConfig:   opts.TlsConfig,

		listenerH:      listenerH,
		clusterManager: clusterManager,
	}

	return proctree.NewProcessGroup(
		proctree.ProcessGroupOptions{
			// Longer init timeout because discovering other nodes
			// could take longer than the default 5 seconds
			InitTimeout: 1 * time.Minute,
			JoinTimeout: 5 * time.Second,
		},
		clusterStateStorage,
		proctree.NewProcess(listenerH),
		clusterManager,
		proctree.NewProcess(serverH),
	)
}
