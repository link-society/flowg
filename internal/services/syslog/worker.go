package syslog

import (
	"log/slog"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/utils/sync"
)

type worker struct {
	logger *slog.Logger

	configStorage  *config.Storage
	pipelineRunner *pipelines.Runner

	state workerState

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
