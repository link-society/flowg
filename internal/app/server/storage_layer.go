package server

import (
	"context"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/utils/proctree"
)

type storageLayer struct {
	authStorage   *auth.Storage
	configStorage *config.Storage
	logStorage    *log.Storage

	process proctree.Process
}

func newStorageLayer(
	authDir string,
	configDir string,
	logDir string,
) *storageLayer {
	var (
		authStorage   = auth.NewStorage(auth.OptDirectory(authDir))
		configStorage = config.NewStorage(config.OptDirectory(configDir))
		logStorage    = log.NewStorage(log.OptDirectory(logDir))
	)

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		authStorage,
		configStorage,
		logStorage,
	)

	return &storageLayer{
		authStorage:   authStorage,
		configStorage: configStorage,
		logStorage:    logStorage,

		process: process,
	}
}

func (l *storageLayer) Start() {
	l.process.Start()
}

func (l *storageLayer) Stop() {
	l.process.Stop()
}

func (l *storageLayer) WaitReady(ctx context.Context) error {
	return l.process.WaitReady(ctx)
}

func (l *storageLayer) Join(ctx context.Context) error {
	return l.process.Join(ctx)
}
