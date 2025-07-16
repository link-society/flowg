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

func (s *ManualClusterFormationStrategy) Join(ctx context.Context, resolver LocalEndpointResolverCallback) ([]*ClusterJoinNode, error) {
	if s.JoinNodeID == "" || s.JoinNodeEndpoint == nil {
		return []*ClusterJoinNode{}, nil
	} else {
		return []*ClusterJoinNode{
			{
				JoinNodeID:       s.JoinNodeID,
				JoinNodeEndpoint: s.JoinNodeEndpoint,
			},
		}, nil
	}
}

func (s *ManualClusterFormationStrategy) Leave(ctx context.Context) error {
	return nil
}
