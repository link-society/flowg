package pipelines

import (
	"context"
)

type pipelineCtxKey string

const (
	workerKey pipelineCtxKey = "worker"
)

// getWorker retrieves the runner worker stashed in the context by the actor when
// it begins handling a logMessage; nodes use it to reach storage and the cache.
func getWorker(ctx context.Context) *worker {
	return ctx.Value(workerKey).(*worker)
}
