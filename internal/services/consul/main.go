package consul

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"time"

	"github.com/hashicorp/go-sockaddr"
	"link-society.com/flowg/internal/cluster"
	"link-society.com/flowg/internal/utils/proctree"
)

type ConsulServiceOptions struct {
	NodeId          string
	ServiceName     string
	ConsulUrl       string
	ClusterJoinNode *cluster.ClusterJoinNode
	MgmtBindAddress string
	MgmtTlsEnabled  bool
}

func NewConsulService(opts *ConsulServiceOptions) proctree.Process {
	consulService := proctree.NewProcess(&procHandler{
		logger: slog.Default().With(
			slog.String("channel", "consul"),
			slog.Group("consul",
				slog.String("consulUrl", opts.ConsulUrl),
			),
		),

		opts: opts,

		LocalEndpointResolver: func() (*url.URL, error) {
			host, port, err := net.SplitHostPort(opts.MgmtBindAddress)
			if err != nil {
				return nil, fmt.Errorf("failed to bind address: %w", err)
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

			var localEndpoint url.URL
			if opts.MgmtTlsEnabled {
				localEndpoint = url.URL{
					Scheme: "https",
					Host:   net.JoinHostPort(host, port),
				}
			} else {
				localEndpoint = url.URL{
					Scheme: "http",
					Host:   net.JoinHostPort(host, port),
				}
			}
			return &localEndpoint, nil
		},
	})

	return proctree.NewProcessGroup(
		proctree.ProcessGroupOptions{
			// Longer init timeout because discovering other nodes
			// could take longer than the default 5 seconds
			InitTimeout: 1 * time.Minute,
			JoinTimeout: 5 * time.Second,
		},
		consulService,
	)

}
