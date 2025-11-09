package pipelines

import (
	"context"
)

type pipelineCtxKey string

const (
	workerKey pipelineCtxKey = "worker"
)

func getWorker(ctx context.Context) *worker {
	return ctx.Value(workerKey).(*worker)
}
