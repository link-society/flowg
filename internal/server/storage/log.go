package storage

import (
	"log/slog"

	"link-society.com/flowg/internal/data/logstorage"
)

type logA struct {
	Storage *logstorage.Storage
	dir     string

	startErrC chan struct{}
	stopErrC  chan struct{}
}

func newLogA(dir string) *logA {
	return &logA{
		Storage: logstorage.NewStorage(logstorage.DefaultStorageOpts().WithDir(dir)),
		dir:     dir,

		startErrC: make(chan struct{}, 1),
		stopErrC:  make(chan struct{}, 1),
	}
}

func (a *logA) Start() {
	defer close(a.startErrC)

	err := a.Storage.Open()
	if err != nil {
		slog.Error(
			"Failed to open logs database",
			"channel", "logstorage",
			"path", a.dir,
			"error", err,
		)
		a.startErrC <- struct{}{}
		return
	}
}

func (a *logA) Stop() {
	defer close(a.stopErrC)

	if a.Storage == nil {
		return
	}

	err := a.Storage.Close()
	if err != nil {
		slog.Error(
			"Failed to close logs database",
			"channel", "logstorage",
			"path", a.dir,
			"error", err,
		)

		a.stopErrC <- struct{}{}
		return
	}
}

func (a *logA) StartErrC() <-chan struct{} {
	return a.startErrC
}

func (a *logA) StopErrC() <-chan struct{} {
	return a.stopErrC
}
