package pipelines

import (
	"context"

	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/data/lognotify"
	"link-society.com/flowg/internal/data/logstorage"
)

type Runner struct {
	ctx context.Context
}

func NewRunner(
	ctx context.Context,
	configStorage *config.Storage,
	logStorage *logstorage.Storage,
	logNotifier *lognotify.LogNotifier,
) *Runner {
	ctx = context.WithValue(ctx, configStorageKey, configStorage)
	ctx = context.WithValue(ctx, logStorageKey, logStorage)
	ctx = context.WithValue(ctx, logNotifierKey, logNotifier)

	return &Runner{
		ctx: ctx,
	}
}

func (r *Runner) Run(pipeline *Pipeline, entry *logstorage.LogEntry) error {
	return pipeline.Process(r.ctx, entry)
}
