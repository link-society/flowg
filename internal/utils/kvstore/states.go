package kvstore

import (
	"github.com/vladopajic/go-actor/actor"

	"github.com/dgraph-io/badger/v4"
)

type workerState interface {
	DoWork(ctx actor.Context, worker *worker) workerState
}

type workerStarting struct {
	dbOpts badger.Options
}

type workerRunning struct {
	db *badger.DB
}

type workerStopping struct {
	db *badger.DB
}

func (s *workerStarting) DoWork(ctx actor.Context, worker *worker) workerState {
	db, err := badger.Open(s.dbOpts)
	if err != nil {
		s.dbOpts.Logger.Errorf("failed to open database: %v", err)
		worker.startCond.Broadcast(ErrStartFailed)
		return nil
	}

	worker.startCond.Broadcast(nil)
	return &workerRunning{db: db}
}

func (s *workerRunning) DoWork(ctx actor.Context, worker *worker) workerState {
	select {
	case <-ctx.Done():
		return &workerStopping{db: s.db}

	case msg, ok := <-worker.mbox.ReceiveC():
		if !ok {
			return &workerStopping{db: s.db}
		}

		go func() {
			err := msg.operation.Handle(s.db)
			msg.replyTo <- err
			close(msg.replyTo)
		}()

		return s
	}
}

func (s *workerStopping) DoWork(ctx actor.Context, worker *worker) workerState {
	if err := s.db.Close(); err != nil {
		s.db.Opts().Logger.Errorf("failed to close database: %v", err)
		worker.stopCond.Broadcast(ErrStopFailed)
	} else {
		worker.stopCond.Broadcast(nil)
	}

	return nil
}
