package cluster

import "context"

type DnsClusterFormationStrategy struct {
	NodeID            string
	ServiceName       string
	DnsServiceAddress string
	DnsDomainName     string
}

// Join implements ClusterFormationStrategy.
func (d *DnsClusterFormationStrategy) Join(ctx context.Context, resolver LocalEndpointResolverCallback) ([]*ClusterJoinNode, error) {
	panic("unimplemented")
}

// Leave implements ClusterFormationStrategy.
func (d *DnsClusterFormationStrategy) Leave(ctx context.Context) error {
	panic("unimplemented")
}
