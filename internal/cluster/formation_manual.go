package cluster

import (
	"context"
	"net/url"
)

type ManualClusterFormationStrategy struct {
	JoinNodeID       string
	JoinNodeEndpoint *url.URL
}

var _ ClusterFormationStrategy = (*ManualClusterFormationStrategy)(nil)

func (s *ManualClusterFormationStrategy) Join(ctx context.Context, resolver LocalEndpointResolverCallback) (*ClusterJoinNode, error) {
	return &ClusterJoinNode{
		JoinNodeID:       s.JoinNodeID,
		JoinNodeEndpoint: s.JoinNodeEndpoint,
	}, nil
}

func (s *ManualClusterFormationStrategy) Leave(ctx context.Context, node *ClusterJoinNode) error {
	return nil
}
