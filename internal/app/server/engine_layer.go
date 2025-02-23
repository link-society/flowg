package server

import (
	"context"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
)

type engineLayer struct {
	logNotifier    *lognotify.LogNotifier
	pipelineRunner *pipelines.Runner

	process proctree.Process
}

func newEngineLayer(storageLayer *storageLayer) *engineLayer {
	logNotifier := lognotify.NewLogNotifier()
	pipelineRunner := pipelines.NewRunner(
		storageLayer.configStorage,
		storageLayer.logStorage,
		logNotifier,
	)

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		logNotifier,
		pipelineRunner,
	)

	return &engineLayer{
		logNotifier:    logNotifier,
		pipelineRunner: pipelineRunner,

		process: process,
	}
}

func (e *engineLayer) Start() {
	e.process.Start()
}

func (e *engineLayer) Stop() {
	e.process.Stop()
}

func (e *engineLayer) WaitReady(ctx context.Context) error {
	return e.process.WaitReady(ctx)
}

func (e *engineLayer) Join(ctx context.Context) error {
	return e.process.Join(ctx)
}
