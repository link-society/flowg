package server

import (
	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
)

type engineLayer struct {
	logNotifier    *lognotify.LogNotifier
	pipelineRunner *pipelines.Runner

	actor actor.Actor
}

func newEngineLayer(storageLayer *storageLayer) *engineLayer {
	logNotifier := lognotify.NewLogNotifier()
	pipelineRunner := pipelines.NewRunner(
		storageLayer.configStorage,
		storageLayer.logStorage,
		logNotifier,
	)

	rootA := actor.Combine(logNotifier, pipelineRunner).
		WithOptions(actor.OptStopTogether()).
		Build()

	return &engineLayer{
		logNotifier:    logNotifier,
		pipelineRunner: pipelineRunner,

		actor: rootA,
	}
}

func (e *engineLayer) Start() {
	e.actor.Start()
}

func (e *engineLayer) Stop() {
	e.actor.Stop()
}
