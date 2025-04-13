package filestore

import (
	"sync"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"
)

type procHandler struct {
	baseDir   string
	inMemory  bool
	extension string

	mbox actor.MailboxReceiver[message]

	mu      sync.Mutex
	backend backend
	cache   *sync.Map
}

var _ proctree.ProcessHandler = (*procHandler)(nil)

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	if h.inMemory {
		h.backend = newMemBackend()
	} else {
		var err error
		h.backend, err = newFsBackend(h.baseDir)
		if err != nil {
			return proctree.Terminate(err)
		}
	}

	h.cache = &sync.Map{}

	keys, err := h.backend.ListFiles()
	if err != nil {
		return proctree.Terminate(err)
	}

	for _, key := range keys {
		content, err := h.backend.ReadFile(key)
		if err != nil {
			return proctree.Terminate(err)
		}

		h.cache.Store(key, content)
	}

	return proctree.Continue()
}

func (h *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	select {
	case <-ctx.Done():
		return proctree.Terminate(ctx.Err())

	case msg, ok := <-h.mbox.ReceiveC():
		if !ok {
			return proctree.Terminate(nil)
		}

		go msg.Handle(h)

		return proctree.Continue()
	}
}

func (h *procHandler) Terminate(ctx actor.Context, err error) error {
	return err
}
