package storage

import (
	"log/slog"

	"link-society.com/flowg/internal/app/bootstrap"

	"link-society.com/flowg/internal/data/config"
)

type configA struct {
	Storage *config.Storage

	startErrC chan struct{}
	stopErrC  chan struct{}
}

func newConfigA(dir string) *configA {
	return &configA{
		Storage: config.NewStorage(config.DefaultStorageOpts().WithDir(dir)),

		startErrC: make(chan struct{}, 1),
		stopErrC:  make(chan struct{}, 1),
	}
}

func (a *configA) Start() {
	defer close(a.startErrC)

	if err := bootstrap.DefaultPipeline(a.Storage); err != nil {
		slog.Error(
			"Failed to bootstrap default pipeline",
			"channel", "config",
			"error", err,
		)
		a.startErrC <- struct{}{}
		return
	}
}

func (a *configA) Stop() {
	close(a.stopErrC)
}

func (a *configA) StartErrC() <-chan struct{} {
	return a.startErrC
}

func (a *configA) StopErrC() <-chan struct{} {
	return a.stopErrC
}
