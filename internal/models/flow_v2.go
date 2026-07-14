package models

// FlowGraphV2 is the current on-disk and API shape of a pipeline's flow graph: a
// set of typed nodes connected by edges. It is what the pipelines engine compiles
// into a runnable pipeline. The minor version distinguishes the original 2.0
// (filter-DSL switch conditions) from 2.1 (expr-lang conditions).
type FlowGraphV2 struct {
	MajorVersion int           `json:"version" default:"2"`
	MinorVersion int           `json:"version.minor" default:"1"`
	HasLayout    bool          `json:"hasLayout" default:"false"`
	Nodes        []*FlowNodeV2 `json:"nodes" required:"true"`
	Edges        []*FlowEdgeV2 `json:"edges" required:"true"`
}

// FlowNodeV2 is a single node of a V2 flow graph. Type selects the behaviour and
// Data carries its type-specific configuration (e.g. the transformer or
// forwarder name).
type FlowNodeV2 struct {
	ID       string            `json:"id" required:"true" minLength:"1"`
	Type     string            `json:"type" required:"true" enum:"source,transform,switch,forwarder,pipeline,router"`
	Position FlowPositionV2    `json:"position" required:"true"`
	Data     map[string]string `json:"data" required:"true"`
}

// FlowEdgeV2 connects a source node to a target node; SourceHandle selects which
// output of the source the edge leaves from.
type FlowEdgeV2 struct {
	ID           string `json:"id" required:"true" minLength:"1"`
	Source       string `json:"source" required:"true" minLength:"1"`
	SourceHandle string `json:"sourceHandle" default:""`
	Target       string `json:"target" required:"true" minLength:"1"`
}

// FlowPositionV2 is a node's canvas position in the editor.
type FlowPositionV2 struct {
	X float64 `json:"x" required:"true"`
	Y float64 `json:"y" required:"true"`
}
