package consul

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"

	"github.com/hashicorp/go-sockaddr"
	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/proctree"
)

type ConsulServiceOptions struct {
	BindAddress     string
	NodeId          string
	ServiceName     string
	ConsulUrl       string
	ClusterJoinNode *models.ClusterJoinNode
	MgmtBindAddress string
}

func NewConsulService(opts *ConsulServiceOptions) proctree.Process {
	return proctree.NewProcess(&procHandler{
		logger: slog.Default().With(
			slog.String("channel", "consul"),
			slog.Group("consul",
				slog.String("consulUrl", opts.ConsulUrl),
			),
		),

		opts: opts,

		LocalEndpointResolver: func() (*url.URL, error) {
			host, port, err := net.SplitHostPort(opts.BindAddress)
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

			localEndpoint := &url.URL{
				Scheme: "http",
				Host:   net.JoinHostPort(host, port),
			}
			return localEndpoint, nil
		},
	})
}
