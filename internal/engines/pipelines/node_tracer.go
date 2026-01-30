package pipelines

import (
	"context"
)

const TRACER_KEY = "tracer_key"

type NodeTrace struct {
	NodeID string            `json:"nodeID"`
	Input  map[string]string `json:"input"`
	Output map[string]string `json:"output,omitempty"`
}

type NodeTracer struct {
	Trace []NodeTrace `json:"trace"`
}

func WithTracer(ctx context.Context, tracer *NodeTracer) context.Context {
	return context.WithValue(ctx, TRACER_KEY, tracer)
}

func GetTracer(ctx context.Context) *NodeTracer {
	m := ctx.Value(TRACER_KEY)
	if m == nil {
		return nil
	}
	return m.(*NodeTracer)
}
