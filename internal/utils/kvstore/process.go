package kvstore

import (
	"errors"

	"github.com/dgraph-io/badger/v4"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"
)

type procHandler struct {
	dbOpts badger.Options
	db     *badger.DB

	mbox actor.MailboxReceiver[message]
}

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	db, err := badger.Open(h.dbOpts)
	if err != nil {
		h.dbOpts.Logger.Errorf("failed to open database: %v", err)
		return proctree.Terminate(err)
	}

	h.db = db

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

		go func() {
			err := msg.operation.Handle(h.db)
			msg.replyTo <- err
			close(msg.replyTo)
		}()

		return proctree.Continue()
	}
}

func (h *procHandler) Terminate(ctx actor.Context, err error) error {
	if !h.db.IsClosed() {
		newError := h.db.Close()
		if newError != nil {
			h.dbOpts.Logger.Errorf("failed to close database: %v", newError)
			err = errors.Join(err, ErrStopFailed)
		}
	}

	return err
}
