package models

import "net/url"

type ClusterJoinNode struct {
	JoinNodeID       string
	JoinNodeEndpoint *url.URL
}
