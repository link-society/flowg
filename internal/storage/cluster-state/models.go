package clusterstate

type NodeState struct {
	NodeID   string                          `json:"node_id"`
	LastSync map[string][]NamespaceSyncState `json:"last_sync"`
}

type NamespaceSyncState struct {
	Namespace string `json:"namespace"`
	Since     uint64 `json:"since"`
}
