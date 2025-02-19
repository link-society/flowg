package mgmt

import (
	"log/slog"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/sync"
)

type worker struct {
	logger *slog.Logger

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
