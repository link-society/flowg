package consul

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/utils/proctree"

	"github.com/hashicorp/consul/api"
)

type procHandler struct {
	consulClient *api.Client
	logger       *slog.Logger
	opts         *ConsulServiceOptions

	LocalEndpointResolver func() (*url.URL, error)
}

const (
	getNodesMaxRetries  = 10
	healthCheckPath     = "/health"
	healthCheckInterval = 5 * time.Second
	healthCheckTimeout  = 1 * time.Second
	shutdownTimeout     = 5 * time.Second
)

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	// If no consul url is provided then stop the consul service as it is not needed
	if h.opts.ConsulUrl == "" {
		h.logger.InfoContext(ctx, "no consul url provided")
		return proctree.Continue()
	}

	// Register node with Consul
	if err := h.registerNode(ctx); err != nil {
		h.logger.ErrorContext(
			ctx,
			"failed to start Consul service",
			slog.Any("error", err),
		)
		return proctree.Terminate(err)
	}

	// Set the JoinNode in ClusterJoinNode
	err := h.setJoinNodes(ctx)
	if err != nil {
		/* Log the error but don't terminate the process
		because the first node that starts up in the cluster
		will never find any other nodes */
		h.logger.WarnContext(
			ctx,
			"failed to get service nodes from consul",
			slog.Any("error", err),
		)
	}

	return proctree.Continue()
}

func (h *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *procHandler) Terminate(ctx actor.Context, err error) error {
	if h.opts.ConsulUrl == "" {
		return err
	}

	h.logger.InfoContext(ctx, "Deregistering service with consul")

	if deregisterErr := h.consulClient.Agent().ServiceDeregister(h.opts.NodeId); err != nil {
		h.logger.ErrorContext(
			ctx,
			"failed to deregister from consul",
			slog.String("error", deregisterErr.Error()),
		)
		err = errors.Join(err, deregisterErr)
	}

	return err
}

func (h *procHandler) registerNode(ctx actor.Context) error {
	localEndpoint, err := h.LocalEndpointResolver()
	if err != nil {
		return err
	}

	// Create a Consul client
	config := api.DefaultConfig()
	config.Address = h.opts.ConsulUrl
	h.consulClient, err = api.NewClient(config)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"failed to create Consul client",
			slog.Any("error", err),
		)
		return err
	}

	var mgmtPortString string
	_, mgmtPortString, err = net.SplitHostPort(h.opts.MgmtBindAddress)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"failed to split manegemnt bind address",
			slog.Any("error", err),
		)
		return err
	}

	var mgmtPort int
	mgmtPort, err = strconv.Atoi(mgmtPortString)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"failed to convert management port from string to int",
			slog.Any("error", err),
		)
		return err
	}

	mgmtEndpoint := h.getManagementEndpointUrl(ctx, localEndpoint.Hostname(), mgmtPortString)

	// Define the service registration
	registration := &api.AgentServiceRegistration{
		ID:      h.opts.NodeId,
		Name:    h.opts.ServiceName,
		Address: localEndpoint.Hostname(),
		Port:    mgmtPort,
		Check: &api.AgentServiceCheck{
			Interval: healthCheckInterval.String(),
			HTTP:     fmt.Sprintf("%s://%s:%d%s", localEndpoint.Scheme, localEndpoint.Hostname(), mgmtPort, healthCheckPath),
			Timeout:  healthCheckTimeout.String(),
		},
		Meta: map[string]string{
			"mgmtEndpoint": mgmtEndpoint,
		},
	}

	// Register the service with Consul
	if err = h.consulClient.Agent().ServiceRegister(registration); err != nil {
		h.logger.ErrorContext(
			ctx,
			"failed to register service with Consul",
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

// setJoinNodes() retries with exponential backoff with jitter to fetch other nodes in the cluster using consul
// and sets one node as a JoinNode in the ClusterJoinNode
// ClusterJoinNode is shared between ConsulService and ManagementServer
func (h *procHandler) setJoinNodes(ctx actor.Context) error {
	retryCount := 0
	delay := 100 * time.Millisecond

	for retryCount <= getNodesMaxRetries {
		serviceEntries, _, err := h.consulClient.Health().Service(h.opts.ServiceName, "", false, nil)
		if err != nil {
			h.logger.ErrorContext(
				ctx,
				"failed to get nodes from consul",
				slog.Any("error", err),
			)
			return err
		}

		for _, serviceEntry := range serviceEntries {
			if serviceEntry.Service.ID != h.opts.NodeId {
				h.opts.ClusterJoinNode.JoinNodeID = serviceEntry.Service.ID

				var mgmtEndpointUrl *url.URL
				mgmtEndpointUrl, err = url.Parse(serviceEntry.Service.Meta["mgmtEndpoint"])

				if err != nil {
					h.logger.ErrorContext(
						ctx,
						"failed to parse management endpoint",
						slog.Any("error", err),
					)
					return err
				}

				h.opts.ClusterJoinNode.JoinNodeEndpoint = &url.URL{
					Scheme: mgmtEndpointUrl.Scheme,
					Host:   net.JoinHostPort(serviceEntry.Service.ID, strconv.Itoa(serviceEntry.Service.Port)),
				}
				return nil
			}
		}

		retryCount++
		if retryCount <= getNodesMaxRetries {
			h.logger.InfoContext(ctx, "did not find other nodes, will try again with a delay")
			time.Sleep(delay)
			// Add jitter to the delay
			delay += time.Duration(rand.IntN(int(delay / 4)))
		}
	}

	return fmt.Errorf("failed to find other nodes")
}

func (h *procHandler) getManagementEndpointUrl(ctx actor.Context, hostname, port string) string {
	var mgmtEndpointUrl url.URL
	if h.opts.MgmtTlsEnabled {
		mgmtEndpointUrl = url.URL{
			Scheme: "https",
			Host:   net.JoinHostPort(hostname, port),
		}
	} else {
		mgmtEndpointUrl = url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(hostname, port),
		}
	}
	return mgmtEndpointUrl.String()
}
