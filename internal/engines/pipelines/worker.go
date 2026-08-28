package pipelines

import (
	"sync"

	"github.com/vladopajic/go-actor/actor"

	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/engines/lognotify"
)

// worker is the runner's actor body. It owns the cache of compiled pipelines and
// the storage handles every node needs; serialising message handling keeps cache
// access safe.
type worker struct {
	mbox actor.MailboxReceiver[message]

	configStorage storage.ConfigStorage
	logStorage    storage.LogStorage
	logNotifier   lognotify.LogNotifier

	cache   map[string]*Pipeline
	cacheMu sync.Mutex
}

var _ actor.Worker = (*worker)(nil)

func (w *worker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case msg, ok := <-w.mbox.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		msg.handle(ctx, w)

		return actor.WorkerContinue
	}
}

// getPipeline retrieves a compiled pipeline from the worker's cache by name.
func (w *worker) getPipeline(_ctx actor.Context, pipelineName string) (*Pipeline, error) {
	w.cacheMu.Lock()
	defer w.cacheMu.Unlock()

	pipeline, exists := w.cache[pipelineName]
	if !exists {
		return nil, &PipelineNotFoundError{Pipeline: pipelineName}
	}

	return pipeline, nil
}

// buildPipeline compiles a new pipeline from storage and adds it to the
// worker's cache.
func (w *worker) buildPipeline(ctx actor.Context, pipelineName string) error {
	pipeline, err := BuildFromStorage(ctx, w.configStorage, pipelineName)
	if err != nil {
		return err
	}
	if pipeline == nil {
		return &PipelineNotFoundError{Pipeline: pipelineName}
	}

	if err := pipeline.Init(ctx); err != nil {
		_ = pipeline.Close(ctx)
		return err
	}

	w.cacheMu.Lock()
	w.cache[pipelineName] = pipeline
	w.cacheMu.Unlock()

	return nil
}

// closePipeline removes a pipeline from the worker's cache and closes it.
func (w *worker) closePipeline(ctx actor.Context, pipelineName string) error {
	w.cacheMu.Lock()
	pipeline, exists := w.cache[pipelineName]
	if !exists {
		w.cacheMu.Unlock()
		return &PipelineNotFoundError{Pipeline: pipelineName}
	}
	delete(w.cache, pipelineName)
	w.cacheMu.Unlock()

	return pipeline.Close(ctx)
}
