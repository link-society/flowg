package kvstore

import (
	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/sync"
)

type workerMbox = actor.MailboxReceiver[message]

type worker struct {
	state workerState
	mbox  workerMbox

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
