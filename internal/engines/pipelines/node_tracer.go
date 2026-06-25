package pipelines

import (
	"context"
	"fmt"

	"link-society.com/flowg/internal/models"
)

const TRACER_KEY = "tracer_key"

// NodeTrace is the record of a single node's execution during a dry run: its
// input fields, the records it emitted, and any error.
type NodeTrace struct {
	NodeID string              `json:"nodeID"`
	Input  map[string]string   `json:"input"`
	Output []map[string]string `json:"output"`
	Error  *string             `json:"error"`
}

// NodeTracer collects the traces of a dry run alongside the flow it ran, so the
// UI can replay how a record travelled through the pipeline.
type NodeTracer struct {
	Flow  models.FlowGraphV2 `json:"flow"`
	Trace []NodeTrace        `json:"trace"`
}

// WithTracer attaches a tracer to the context, switching processing into dry-run
// mode.
func WithTracer(ctx context.Context, tracer *NodeTracer) context.Context {
	return context.WithValue(ctx, TRACER_KEY, tracer)
}

// GetTracer returns the tracer carried by the context, or nil during a normal
// (non-dry) run.
func GetTracer(ctx context.Context) *NodeTracer {
	m := ctx.Value(TRACER_KEY)
	if m == nil {
		return nil
	}
	return m.(*NodeTracer)
}

// TraceError renders an error as a pointer to its string, or nil when there is
// no error, matching the JSON shape of NodeTrace.Error.
func TraceError(err error) *string {
	if err == nil {
		return nil
	}

	errMsg := fmt.Sprint(err)
	return &errMsg
}
