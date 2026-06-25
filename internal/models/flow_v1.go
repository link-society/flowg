package models

// FlowGraphV1 is the legacy (version 1) on-disk shape of a pipeline's flow graph.
// It is kept only so old pipelines can be read and upgraded to V2 on load; new
// pipelines are never written in this format. See flow_convert.go.
type FlowGraphV1 struct {
	Version int           `json:"version" default:"1"`
	Nodes   []*FlowNodeV1 `json:"nodes"`
	Edges   []*FlowEdgeV1 `json:"edges"`
}

// FlowNodeV1 is a node in a V1 flow graph.
type FlowNodeV1 struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Position FlowPositionV1    `json:"position"`
	Data     map[string]string `json:"data"`
}

// FlowEdgeV1 connects two V1 nodes.
type FlowEdgeV1 struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle"`
	Target       string `json:"target"`
}

// FlowPositionV1 is a node's canvas position in the editor.
type FlowPositionV1 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
