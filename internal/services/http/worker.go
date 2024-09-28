package http

import (
	"log/slog"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/sync"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type worker struct {
	logger *slog.Logger

	authStorage   *auth.Storage
	configStorage *config.Storage
	logStorage    *log.Storage

	logNotifier    *lognotify.LogNotifier
	pipelineRunner *pipelines.Runner

	state     workerState
	startCond *sync.CondValue[error]
	stopCond  *sync.CondValue[error]
}

func (w *worker) DoWork(ctx actor.Context) actor.WorkerStatus {
	if w.state == nil {
		return actor.WorkerEnd
	}

	w.state = w.state.DoWork(ctx, w)
	return actor.WorkerContinue
}
