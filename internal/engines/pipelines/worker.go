package pipelines

import (
	"context"
	"log/slog"
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

// acquirePipeline returns the cached build with an in-flight reservation so a
// concurrent retire cannot close it mid-use; callers must call release.
func (w *worker) acquirePipeline(pipelineName string) (*Pipeline, func(), error) {
	w.cacheMu.Lock()
	pipeline, exists := w.cache[pipelineName]
	w.cacheMu.Unlock()

	if !exists || !pipeline.acquire() {
		return nil, nil, &PipelineNotFoundError{Pipeline: pipelineName}
	}

	return pipeline, pipeline.release, nil
}

// retirePipeline marks a build closed and closes it in the background once its
// in-flight uses drained; the returned channel reports the close error.
func (w *worker) retirePipeline(pipeline *Pipeline) <-chan error {
	pipeline.stateMu.Lock()
	pipeline.closed = true
	pipeline.stateMu.Unlock()

	done := make(chan error, 1)
	go func() {
		pipeline.inflight.Wait()
		done <- pipeline.Close(context.Background())
		close(done)
	}()

	return done
}

// buildPipeline compiles a new build from storage, swaps it into the worker's
// cache and retires the previous build, if any. In-flight records finish on
// the old build while new ones pick up the new one.
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
	previous := w.cache[pipelineName]
	w.cache[pipelineName] = pipeline
	w.cacheMu.Unlock()

	if previous != nil {
		done := w.retirePipeline(previous)
		go func() {
			if err := <-done; err != nil {
				slog.ErrorContext(
					ctx,
					"failed to close previous pipeline build",
					"channel", "pipelines",
					"pipeline", pipelineName,
					"error", err.Error(),
				)
			}
		}()
	}

	return nil
}

// closePipeline removes a build from the worker's cache and retires it; the
// returned channel reports the close error once in-flight work drained.
func (w *worker) closePipeline(pipelineName string) (<-chan error, error) {
	w.cacheMu.Lock()
	pipeline, exists := w.cache[pipelineName]
	delete(w.cache, pipelineName)
	w.cacheMu.Unlock()

	if !exists {
		return nil, &PipelineNotFoundError{Pipeline: pipelineName}
	}

	return w.retirePipeline(pipeline), nil
}
