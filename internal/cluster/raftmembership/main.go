package raftmembership

import (
	"context"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"github.com/hashicorp/raft"
)

type Server struct {
	mbox    actor.MailboxSender[*ChangeRequest]
	process proctree.Process
}

func NewServer(raft *raft.Raft) *Server {
	mbox := actor.NewMailbox[*ChangeRequest]()
	handler := &procHandler{raft: raft, mbox: mbox}

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewActorProcess(mbox),
		proctree.NewProcess(handler),
	)

	return &Server{mbox: mbox, process: process}
}

func (s *Server) Start() {
	s.process.Start()
}

func (s *Server) Stop() {
	s.process.Stop()
}

func (s *Server) WaitReady(ctx context.Context) error {
	return s.process.WaitReady(ctx)
}

func (s *Server) Join(ctx context.Context) error {
	return s.process.Join(ctx)
}

func (s *Server) SendRequest(ctx context.Context, req *ChangeRequest) error {
	err := s.mbox.Send(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()

	case err := <-req.done:
		return err
	}
}
