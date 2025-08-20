package pipelines

import (
	"context"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type pipelineCtxKey string

const (
	configStorageKey pipelineCtxKey = "configStorage"
	logStorageKey    pipelineCtxKey = "logStorage"
	logNotifierKey   pipelineCtxKey = "logNotifier"
)

func getConfigStorage(ctx context.Context) config.Storage {
	return ctx.Value(configStorageKey).(config.Storage)
}

func getLogStorage(ctx context.Context) log.Storage {
	return ctx.Value(logStorageKey).(log.Storage)
}

func getLogNotifier(ctx context.Context) lognotify.LogNotifier {
	return ctx.Value(logNotifierKey).(lognotify.LogNotifier)
}
