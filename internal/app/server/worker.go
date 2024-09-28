package server

import (
	"log/slog"

	"github.com/vladopajic/go-actor/actor"
)

type worker struct {
	state  workerState
	logger *slog.Logger

	storageLayer *storageLayer
	engineLayer  *engineLayer
	serviceLayer *serviceLayer

	failure bool
}

func (w *worker) DoWork(ctx actor.Context) actor.WorkerStatus {
	if w.state == nil {
		return actor.WorkerEnd
	}

	w.state = w.state.DoWork(ctx, w)
	return actor.WorkerContinue
}
