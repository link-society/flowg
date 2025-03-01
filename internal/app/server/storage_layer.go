package server

import (
	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/utils/proctree"
)

type storageLayer struct {
	proctree.Process

	authStorage   *auth.Storage
	configStorage *config.Storage
	logStorage    *log.Storage
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
		Process: process,

		authStorage:   authStorage,
		configStorage: configStorage,
		logStorage:    logStorage,
	}
}
