package pipelines

import (
	"context"

	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/data/lognotify"
	"link-society.com/flowg/internal/data/logstorage"
)

func Run(
	pipeline *Pipeline,
	ctx context.Context,
	configStorage *config.Storage,
	logStorage *logstorage.Storage,
	logNotifier *lognotify.LogNotifier,
	entry *logstorage.LogEntry,
) error {
	ctx = context.WithValue(ctx, configStorageKey, configStorage)
	ctx = context.WithValue(ctx, logStorageKey, logStorage)
	ctx = context.WithValue(ctx, logNotifierKey, logNotifier)

	return pipeline.Process(ctx, entry)
}
