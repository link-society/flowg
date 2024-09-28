package filestore

import (
	"sync"

	"github.com/vladopajic/go-actor/actor"
)

type workerState interface {
	DoWork(ctx actor.Context, worker *worker) workerState
}

type workerStarting struct {
	baseDir   string
	inMemory  bool
	extension string
}

type workerRunning struct {
	mu      sync.Mutex
	backend backend
	cache   *sync.Map

	extension string
}

type workerStopping struct{}

func (s *workerStarting) DoWork(ctx actor.Context, worker *worker) workerState {
	var backend backend

	if s.inMemory {
		backend = newMemBackend()
	} else {
		var err error
		backend, err = newFsBackend(s.baseDir)
		if err != nil {
			worker.startCond.Broadcast(err)
			return &workerStopping{}
		}
	}

	cache := &sync.Map{}

	keys, err := backend.ListFiles()
	if err != nil {
		worker.startCond.Broadcast(err)
		return &workerStopping{}
	}

	for _, key := range keys {
		content, err := backend.ReadFile(key)
		if err != nil {
			worker.startCond.Broadcast(err)
			return &workerStopping{}
		}

		cache.Store(key, content)
	}

	worker.startCond.Broadcast(nil)
	return &workerRunning{
		mu:      sync.Mutex{},
		backend: backend,
		cache:   cache,

		extension: s.extension,
	}
}

func (s *workerRunning) DoWork(ctx actor.Context, worker *worker) workerState {
	select {
	case <-ctx.Done():
		return &workerStopping{}

	case msg, ok := <-worker.mbox.ReceiveC():
		if !ok {
			return &workerStopping{}
		}

		go msg.Handle(s)

		return s
	}
}

func (s *workerStopping) DoWork(ctx actor.Context, worker *worker) workerState {
	worker.stopCond.Broadcast(nil)
	return nil
}
