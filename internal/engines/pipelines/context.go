package pipelines

import (
	"context"
)

type pipelineCtxKey string

const (
	workerKey   pipelineCtxKey = "worker"
	pipelineKey pipelineCtxKey = "pipeline"
)

// getWorker retrieves the runner worker stashed in the context by the actor when
// it begins handling a logMessage; nodes use it to reach storage and the cache.
func getWorker(ctx context.Context) *worker {
	return ctx.Value(workerKey).(*worker)
}

// getPipeline retrieves the pipeline stashed in the context, so that nodes can
// access the pipeline's they are part of.
func getPipeline(ctx context.Context) *Pipeline {
	return ctx.Value(pipelineKey).(*Pipeline)
}
