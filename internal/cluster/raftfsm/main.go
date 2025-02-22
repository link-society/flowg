package raftfsm

import (
	"io"

	"github.com/hashicorp/raft"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type fsm struct {
	authStorage   *auth.Storage
	configStorage *config.Storage
	logStorage    *log.Storage
}

func New(
	authStorage *auth.Storage,
	configStorage *config.Storage,
	logStorage *log.Storage,
) raft.FSM {
	return &fsm{
		authStorage:   authStorage,
		configStorage: configStorage,
		logStorage:    logStorage,
	}
}

func (f *fsm) Apply(l *raft.Log) interface{} {
	return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (f *fsm) Restore(r io.ReadCloser) error {
	return nil
}
