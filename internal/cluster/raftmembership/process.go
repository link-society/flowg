package raftmembership

import (
	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/utils/proctree"

	"github.com/hashicorp/raft"
)

type procHandler struct {
	raft *raft.Raft
	mbox actor.MailboxReceiver[*ChangeRequest]
}

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	return proctree.Continue()
}

func (h *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	select {
	case <-ctx.Done():
		return proctree.Terminate(ctx.Err())

	case req, ok := <-h.mbox.ReceiveC():
		if !ok {
			return proctree.Terminate(nil)
		}

		req.notifyDone(req.kind.Handle(h.raft, req.id))

		return proctree.Continue()
	}
}

func (h *procHandler) Terminate(ctx actor.Context, err error) error {
	return err
}
