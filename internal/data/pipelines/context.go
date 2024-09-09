package pipelines

import (
	"context"

	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/data/lognotify"
	"link-society.com/flowg/internal/data/logstorage"
)

type pipelineCtxKey string

const (
	configStorageKey pipelineCtxKey = "configStorage"
	logStorageKey    pipelineCtxKey = "logStorage"
	logNotifierKey   pipelineCtxKey = "logNotifier"
)

func getTransformerSystem(ctx context.Context) *config.TransformerSystem {
	configStorage := ctx.Value(configStorageKey).(*config.Storage)
	return config.NewTransformerSystem(configStorage)
}

func getPipelineSystem(ctx context.Context) *config.PipelineSystem {
	configStorage := ctx.Value(configStorageKey).(*config.Storage)
	return config.NewPipelineSystem(configStorage)
}

func getCollectorSystem(ctx context.Context) *logstorage.CollectorSystem {
	logStorage := ctx.Value(logStorageKey).(*logstorage.Storage)
	return logstorage.NewCollectorSystem(logStorage)
}

func getLogNotifier(ctx context.Context) *lognotify.LogNotifier {
	return ctx.Value(logNotifierKey).(*lognotify.LogNotifier)
}
