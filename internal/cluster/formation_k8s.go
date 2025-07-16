package cluster

import (
	"context"
	"log/slog"

	"fmt"
	"strconv"

	"net"
	"net/url"

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

func (s *KubernetesClusterFormationStrategy) Join(ctx context.Context, resolver LocalEndpointResolverCallback) ([]*ClusterJoinNode, error) {
	logger := slog.Default().With(slog.String("channel", "cluster.formation.k8s"))

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

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

	joinNodes := []*ClusterJoinNode{}

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
					joinNodes = append(joinNodes, &ClusterJoinNode{
						JoinNodeID: endpoint.TargetRef.Name,
						JoinNodeEndpoint: &url.URL{
							Scheme: scheme,
							Host:   net.JoinHostPort(addr, strconv.Itoa(int(*port.Port))),
						},
					})
				}
			}
		}
	}

	return joinNodes, nil
}

func (s *KubernetesClusterFormationStrategy) Leave(ctx context.Context) error {
	return nil
}
