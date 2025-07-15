package cluster

import (
	"context"
	"log/slog"
	"time"

	"fmt"
	"strconv"

	"net"
	"net/url"

	"github.com/avast/retry-go/v4"

	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesClusterFormationStrategy struct {
	NodeID           string
	ServiceNamespace string
	ServiceName      string
	ServicePortName  string
}

var _ ClusterFormationStrategy = (*KubernetesClusterFormationStrategy)(nil)

func (s *KubernetesClusterFormationStrategy) Join(ctx context.Context, resolver LocalEndpointResolverCallback) (*ClusterJoinNode, error) {
	const (
		getNodesMaxRetries = 5
		getNodesDelay      = 5 * time.Second
		getNodesMaxJitter  = getNodesDelay / 4
	)

	logger := slog.Default().With(slog.String("channel", "cluster.k8s"))

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	joinNode, err := retry.DoWithData(
		func() (*ClusterJoinNode, error) {
			logger.InfoContext(ctx, "discover available nodes from Kubernetes")
			endpointSlices, err := clientset.
				DiscoveryV1().
				EndpointSlices(s.ServiceNamespace).
				List(ctx, metav1.ListOptions{
					LabelSelector: fmt.Sprintf("kubernetes.io/service-name=%s", s.ServiceName),
				})
			if err != nil {
				return nil, fmt.Errorf("failed to list endpoints: %w", err)
			}

			if len(endpointSlices.Items) == 0 {
				return nil, fmt.Errorf(
					"no endpoints found for service %s in namespace %s",
					s.ServiceName,
					s.ServiceNamespace,
				)
			}

			type endpointInfo struct {
				nodeID  string
				address string
				port    int32
				scheme  string
			}

			var endpoints []endpointInfo

			for _, slice := range endpointSlices.Items {
				var port discoveryv1.EndpointPort
				portFound := false

				for _, p := range slice.Ports {
					if p.Name != nil && p.Port != nil && *p.Name == s.ServicePortName {
						port = p
						portFound = true
						break
					}
				}

				if !portFound {
					return nil, fmt.Errorf(
						"no port found with name %s in service %s in namespace %s",
						s.ServicePortName,
						s.ServiceName,
						s.ServiceNamespace,
					)
				}

				scheme := "http"
				if port.AppProtocol != nil {
					if *port.AppProtocol != "https" && *port.AppProtocol != "http" {
						return nil, fmt.Errorf(
							"unsupported protocol %s for port %s in service %s in namespace %s",
							*port.AppProtocol,
							s.ServicePortName,
							s.ServiceName,
							s.ServiceNamespace,
						)
					}
					scheme = *port.AppProtocol
				}

				for _, endpoint := range slice.Endpoints {
					if endpoint.Conditions.Ready != nil && !*endpoint.Conditions.Ready {
						continue // Skip endpoints that are not ready
					}

					if endpoint.TargetRef.Name != s.NodeID {
						for _, addr := range endpoint.Addresses {
							endpoints = append(endpoints, endpointInfo{
								nodeID:  endpoint.TargetRef.Name,
								address: addr,
								port:    *port.Port,
								scheme:  scheme,
							})
						}
					}
				}
			}

			if len(endpoints) == 0 {
				return nil, fmt.Errorf(
					"no valid endpoint found for service %s in namespace %s",
					s.ServiceName,
					s.ServiceNamespace,
				)
			}

			logger.InfoContext(
				ctx,
				"found join node in Kubernetes",
				slog.String("node.id", endpoints[0].nodeID),
				slog.String("endpoint", fmt.Sprintf("%s://%s:%d", endpoints[0].scheme, endpoints[0].address, endpoints[0].port)),
			)

			return &ClusterJoinNode{
				JoinNodeID: endpoints[0].nodeID,
				JoinNodeEndpoint: &url.URL{
					Scheme: endpoints[0].scheme,
					Host:   net.JoinHostPort(endpoints[0].address, strconv.Itoa(int(endpoints[0].port))),
				},
			}, nil
		},
		retry.Attempts(getNodesMaxRetries),
		retry.Delay(getNodesDelay),
		retry.MaxJitter(getNodesMaxJitter),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			delay := retry.FixedDelay(n, err, config)
			delay += retry.RandomDelay(n, err, config)
			return delay
		}),
		retry.OnRetry(func(n uint, err error) {
			logger.WarnContext(
				ctx,
				"retrying to discover nodes from Kubernetes",
				slog.Uint64("attempt", uint64(n)),
				slog.String("error", err.Error()),
			)
		}),
	)

	if err != nil {
		logger.WarnContext(
			ctx,
			"failed to get join nodes from Kubernetes",
			slog.String("error", err.Error()),
		)
	}

	if joinNode == nil {
		joinNode = &ClusterJoinNode{}
	}

	return joinNode, nil
}

func (s *KubernetesClusterFormationStrategy) Leave(ctx context.Context, node *ClusterJoinNode) error {
	return nil
}
