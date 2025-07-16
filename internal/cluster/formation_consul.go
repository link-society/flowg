package cluster

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"time"

	"net"
	"net/url"

	"github.com/hashicorp/consul/api"
)

type ConsulClusterFormationStrategy struct {
	NodeID         string
	ServiceName    string
	ServiceAddress string
	ServiceTls     bool
	ConsulUrl      string

	client *api.Client
}

var _ ClusterFormationStrategy = (*ConsulClusterFormationStrategy)(nil)

func (s *ConsulClusterFormationStrategy) Join(ctx context.Context, resolver LocalEndpointResolverCallback) ([]*ClusterJoinNode, error) {
	const (
		healthCheckInterval = 5 * time.Second
		healthCheckTimeout  = 1 * time.Second
	)

	logger := slog.Default().With(slog.String("channel", "cluster.formation.consul"))

	localEndpoint, err := resolver()
	if err != nil {
		return nil, err
	}

	config := api.DefaultConfig()
	config.Address = s.ConsulUrl
	s.client, err = api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	_, svcPort, err := net.SplitHostPort(s.ServiceAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service address: %w", err)
	}

	svcPortNumber, err := strconv.Atoi(svcPort)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service port: %w", err)
	}

	var serviceEndpoint *url.URL
	if s.ServiceTls {
		serviceEndpoint = &url.URL{
			Scheme: "https",
			Host:   net.JoinHostPort(localEndpoint.Hostname(), svcPort),
		}
	} else {
		serviceEndpoint = &url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(localEndpoint.Hostname(), svcPort),
		}
	}

	registration := &api.AgentServiceRegistration{
		ID:      s.NodeID,
		Name:    s.ServiceName,
		Address: localEndpoint.Hostname(),
		Port:    svcPortNumber,
		Check: &api.AgentServiceCheck{
			Interval: healthCheckInterval.String(),
			HTTP:     fmt.Sprintf("%s/health", serviceEndpoint.String()),
			Timeout:  healthCheckTimeout.String(),
		},
		Meta: map[string]string{
			"endpoint": serviceEndpoint.String(),
		},
	}

	logger.InfoContext(ctx, "registering service with Consul")
	if err := s.client.Agent().ServiceRegister(registration); err != nil {
		return nil, fmt.Errorf("failed to register service with Consul: %w", err)
	}

	logger.InfoContext(ctx, "discover available nodes from Consul")
	entries, _, err := s.client.Health().Service(s.ServiceName, "", false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes from Consul: %w", err)
	}

	joinNodes := []*ClusterJoinNode{}

	for _, entry := range entries {
		if entry.Service.ID != s.NodeID {
			endpoint, err := url.Parse(entry.Service.Meta["endpoint"])
			if err != nil {
				return nil, fmt.Errorf("failed to parse service endpoint URL: %w", err)
			}

			logger.InfoContext(
				ctx,
				"found join node in Consul",
				slog.String("node_id", entry.Service.ID),
				slog.String("endpoint", endpoint.String()),
			)
			joinNodes = append(joinNodes, &ClusterJoinNode{
				JoinNodeID:       entry.Service.ID,
				JoinNodeEndpoint: endpoint,
			})
		}
	}

	return joinNodes, nil
}

func (s *ConsulClusterFormationStrategy) Leave(ctx context.Context) error {
	if s.client != nil {
		logger := slog.Default().With(slog.String("channel", "cluster.consul"))

		logger.InfoContext(ctx, "deregistering service from Consul")
		err := s.client.Agent().ServiceDeregister(s.NodeID)
		if err != nil {
			return fmt.Errorf("failed to deregister service from Consul: %w", err)
		}
	}

	return nil
}
