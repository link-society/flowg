package server

import (
	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
)

type engineLayer struct {
	proctree.Process

	logNotifier    *lognotify.LogNotifier
	pipelineRunner *pipelines.Runner
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
		Process: process,

		logNotifier:    logNotifier,
		pipelineRunner: pipelineRunner,
	}
}
