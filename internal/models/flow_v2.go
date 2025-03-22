package models

type FlowGraphV2 struct {
	Version int           `json:"version" default:"2"`
	Nodes   []*FlowNodeV2 `json:"nodes"`
	Edges   []*FlowEdgeV2 `json:"edges"`
}

type FlowNodeV2 struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Position FlowPositionV2    `json:"position"`
	Data     map[string]string `json:"data"`
}

type FlowEdgeV2 struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle"`
	Target       string `json:"target"`
}

type FlowPositionV2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
