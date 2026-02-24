package pipelines

import (
	"context"
	"fmt"

	"link-society.com/flowg/internal/models"
)

const TRACER_KEY = "tracer_key"

type NodeTrace struct {
	NodeID string            `json:"nodeID"`
	Input  map[string]string `json:"input"`
	Output map[string]string `json:"output"`
	Error  *string           `json:"error"`
}

type NodeTracer struct {
	Flow  models.FlowGraphV2 `json:"flow"`
	Trace []NodeTrace        `json:"trace"`
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

func TraceError(err error) *string {
	if err == nil {
		return nil
	}

	errMsg := fmt.Sprint(err)
	return &errMsg
}
