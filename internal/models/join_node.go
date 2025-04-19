package models

import "net/url"

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
