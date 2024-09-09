package config

type FlowGraph struct {
	Nodes []*FlowNode `json:"nodes"`
	Edges []*FlowEdge `json:"edges"`
}

type FlowNode struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Position FlowPosition      `json:"position"`
	Data     map[string]string `json:"data"`
}

type FlowEdge struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle"`
	Target       string `json:"target"`
}

type FlowPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
