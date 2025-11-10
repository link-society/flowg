package models

type FlowGraphV2 struct {
	Version int           `json:"version" default:"2"`
	Nodes   []*FlowNodeV2 `json:"nodes" required:"true"`
	Edges   []*FlowEdgeV2 `json:"edges" required:"true"`
}

type FlowNodeV2 struct {
	ID       string            `json:"id" required:"true"`
	Type     string            `json:"type" required:"true"`
	Position FlowPositionV2    `json:"position" required:"true"`
	Data     map[string]string `json:"data" required:"true"`
}

type FlowEdgeV2 struct {
	ID           string `json:"id" required:"true"`
	Source       string `json:"source" required:"true"`
	SourceHandle string `json:"sourceHandle" default:""`
	Target       string `json:"target" required:"true"`
}

type FlowPositionV2 struct {
	X float64 `json:"x" required:"true"`
	Y float64 `json:"y" required:"true"`
}
