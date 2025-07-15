package cluster

import (
	"fmt"
	"net/url"
)

type ClusterJoinNode struct {
	JoinNodeID       string
	JoinNodeEndpoint *url.URL
}

func NewClusterJoinNode(isAutomaticClusterformation bool, defaultJoinNodeId string, defaultJoinNodeEndpoint *url.URL) *ClusterJoinNode {
	if isAutomaticClusterformation {
		return &ClusterJoinNode{}
	}
	return &ClusterJoinNode{
		JoinNodeID:       defaultJoinNodeId,
		JoinNodeEndpoint: defaultJoinNodeEndpoint,
	}
}

func (c *ClusterJoinNode) IsEmpty() bool {
	return c.JoinNodeID == "" || c.JoinNodeEndpoint == nil
}

func (c *ClusterJoinNode) Address() string {
	return fmt.Sprintf("%s/%s", c.JoinNodeID, c.JoinNodeEndpoint.Host)
}
