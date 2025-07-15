package cluster

import (
	"context"
	"net/url"
)

type LocalEndpointResolverCallback func() (*url.URL, error)

type ClusterFormationStrategy interface {
	Join(ctx context.Context, resolver LocalEndpointResolverCallback) (*ClusterJoinNode, error)
	Leave(ctx context.Context, node *ClusterJoinNode) error
}
