package models

type FlowGraphV1 struct {
	Version int           `json:"version"`
	Nodes   []*FlowNodeV1 `json:"nodes"`
	Edges   []*FlowEdgeV1 `json:"edges"`
}

type FlowNodeV1 struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Position FlowPositionV1    `json:"position"`
	Data     map[string]string `json:"data"`
}

type FlowEdgeV1 struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle"`
	Target       string `json:"target"`
}

type FlowPositionV1 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
