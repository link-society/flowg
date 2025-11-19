package models

type FlowGraphV2 struct {
	MajorVersion int           `json:"version" default:"2"`
	MinorVersion int           `json:"version.minor" default:"1"`
	Nodes        []*FlowNodeV2 `json:"nodes" required:"true"`
	Edges        []*FlowEdgeV2 `json:"edges" required:"true"`
}

type FlowNodeV2 struct {
	ID       string            `json:"id" required:"true" minLength:"1"`
	Type     string            `json:"type" required:"true" enum:"source,transform,switch,forwarder,pipeline,router"`
	Position FlowPositionV2    `json:"position" required:"true"`
	Data     map[string]string `json:"data" required:"true"`
}

type FlowEdgeV2 struct {
	ID           string `json:"id" required:"true" minLength:"1"`
	Source       string `json:"source" required:"true" minLength:"1"`
	SourceHandle string `json:"sourceHandle" default:""`
	Target       string `json:"target" required:"true" minLength:"1"`
}

type FlowPositionV2 struct {
	X float64 `json:"x" required:"true"`
	Y float64 `json:"y" required:"true"`
}
