package server

import (
	"errors"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type storageLayer struct {
	authStorage   *auth.Storage
	configStorage *config.Storage
	logStorage    *log.Storage

	actor actor.Actor
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

	rootA := actor.Combine(authStorage, configStorage, logStorage).
		WithOptions(actor.OptStopTogether()).
		Build()

	return &storageLayer{
		authStorage:   authStorage,
		configStorage: configStorage,
		logStorage:    logStorage,

		actor: rootA,
	}
}

func (a *storageLayer) Start() {
	a.actor.Start()
}

func (a *storageLayer) WaitStarted() error {
	errs := []error{}

	if err := a.authStorage.WaitStarted(); err != nil {
		errs = append(errs, err)
	}

	if err := a.configStorage.WaitStarted(); err != nil {
		errs = append(errs, err)
	}

	if err := a.logStorage.WaitStarted(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (a *storageLayer) Stop() {
	a.actor.Stop()
}

func (a *storageLayer) WaitStopped() error {
	errs := []error{}

	if err := a.logStorage.WaitStopped(); err != nil {
		errs = append(errs, err)
	}

	if err := a.configStorage.WaitStopped(); err != nil {
		errs = append(errs, err)
	}

	if err := a.authStorage.WaitStopped(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
