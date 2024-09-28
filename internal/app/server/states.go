package server

import (
	"log/slog"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/app/bootstrap"
)

type workerState interface {
	DoWork(ctx actor.Context, w *worker) workerState
}

type workerStartingStorageLayer struct{}
type workerStartingEngineLayer struct{}
type workerStartingServiceLayer struct{}
type workerRunning struct{}
type workerStoppingServiceLayer struct{}
type workerStoppingEngineLayer struct{}
type workerStoppingStorageLayer struct{}

func (s *workerStartingStorageLayer) DoWork(ctx actor.Context, w *worker) workerState {
	w.storageLayer.Start()
	err := w.storageLayer.WaitStarted()
	if err != nil {
		w.failure = true
		return nil
	}

	err = bootstrap.DefaultRolesAndUsers(ctx, w.storageLayer.authStorage)
	if err != nil {
		w.logger.ErrorContext(
			ctx,
			"Failed to bootstrap default roles and users",
			slog.String("error", err.Error()),
		)
		w.failure = true
		return &workerStoppingStorageLayer{}
	}

	err = bootstrap.DefaultPipeline(ctx, w.storageLayer.configStorage)
	if err != nil {
		w.logger.ErrorContext(
			ctx,
			"Failed to bootstrap default pipeline",
			slog.String("error", err.Error()),
		)
		w.failure = true
		return &workerStoppingStorageLayer{}
	}

	return &workerStartingEngineLayer{}
}

func (s *workerStartingEngineLayer) DoWork(ctx actor.Context, w *worker) workerState {
	w.engineLayer.Start()
	return &workerStartingServiceLayer{}
}

func (s *workerStartingServiceLayer) DoWork(ctx actor.Context, w *worker) workerState {
	w.serviceLayer.Start()
	err := w.serviceLayer.WaitStarted()
	if err != nil {
		w.failure = true
		return &workerStoppingEngineLayer{}
	}

	return &workerRunning{}
}

func (s *workerRunning) DoWork(ctx actor.Context, w *worker) workerState {
	<-ctx.Done()
	return &workerStoppingServiceLayer{}
}

func (s *workerStoppingServiceLayer) DoWork(ctx actor.Context, w *worker) workerState {
	w.serviceLayer.Stop()
	err := w.serviceLayer.WaitStopped()
	if err != nil {
		w.failure = true
	}

	return &workerStoppingEngineLayer{}
}

func (s *workerStoppingEngineLayer) DoWork(ctx actor.Context, w *worker) workerState {
	w.engineLayer.Stop()
	return &workerStoppingStorageLayer{}
}

func (s *workerStoppingStorageLayer) DoWork(ctx actor.Context, w *worker) workerState {
	w.storageLayer.Stop()
	err := w.storageLayer.WaitStopped()
	if err != nil {
		w.failure = true
	}

	return nil
}
