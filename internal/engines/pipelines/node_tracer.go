package pipelines

import (
	"context"
	"fmt"
	"sync"

	"go.link-society.com/flowg/pkg/models"
)

const TRACER_KEY = "tracer_key"

// NodeTracer collects the traces of a dry run alongside the flow it ran, so the
// UI can replay how a record travelled through the pipeline.
//
// Successors of a node are processed concurrently, so entries are added under a
// lock (see [NodeTracer.add]).
type NodeTracer struct {
	Flow  models.FlowGraphV2         `json:"flow"`
	Trace []models.PipelineNodeTrace `json:"trace"`

	mu sync.Mutex
}

// add appends a trace entry, safely for concurrent callers.
func (t *NodeTracer) add(entry models.PipelineNodeTrace) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Trace = append(t.Trace, entry)
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
// no error, matching the JSON shape of PipelineNodeTrace.Error.
func TraceError(err error) *string {
	if err == nil {
		return nil
	}

	errMsg := fmt.Sprint(err)
	return &errMsg
}
