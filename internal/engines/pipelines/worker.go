package pipelines

import (
	"context"
	"sync"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage"

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

// getOrBuildPipeline returns the cached build for a pipeline, compiling and
// initialising it from storage on first use. The build is closed again if
// initialisation fails.
func (w *worker) getOrBuildPipeline(ctx context.Context, pipelineName string) (*Pipeline, error) {
	w.cacheMu.Lock()
	defer w.cacheMu.Unlock()

	if pipeline, exists := w.cache[pipelineName]; exists {
		return pipeline, nil
	}

	pipeline, err := BuildFromStorage(ctx, w.configStorage, pipelineName)
	if err != nil {
		return nil, err
	}

	if err := pipeline.Init(ctx); err != nil {
		_ = pipeline.Close(ctx)
		return nil, err
	}

	w.cache[pipelineName] = pipeline
	return pipeline, nil
}
