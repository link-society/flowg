package pipelines

import (
	"context"
	"errors"

	"link-society.com/flowg/internal/models"
)

const (
	// DIRECT_ENTRYPOINT is the default source type, used when a record is
	// submitted directly (e.g. via the ingestion API or a parent pipeline).
	DIRECT_ENTRYPOINT = "direct"
	// SYSLOG_ENTRYPOINT is the source type for records arriving from the syslog
	// service.
	SYSLOG_ENTRYPOINT = "syslog"
)

// message is a request handled by the runner actor; each variant knows how to
// service itself against the worker.
type message interface {
	handle(ctx context.Context, w *worker)
}

// logMessage requests that a record be processed by a pipeline. When tracer is
// set the run is a dry run that rebuilds the pipeline from the traced flow and
// records per-node traces.
type logMessage struct {
	replyTo chan<- error

	pipelineName string
	entrypoint   string
	record       *models.LogRecord
	tracer       *NodeTracer
}

// invalidateCacheMessage requests eviction of a single pipeline's cached build.
type invalidateCacheMessage struct {
	replyTo      chan<- error
	pipelineName string
}

// invalidateAllCacheMessage requests eviction of every cached pipeline build.
type invalidateAllCacheMessage struct {
	replyTo chan<- error
}

func (msg logMessage) handle(ctx context.Context, w *worker) {
	go func() {
		var pipeline *Pipeline
		var err error
		defer close(msg.replyTo)

		ctx := context.WithValue(ctx, workerKey, w)
		if msg.tracer != nil {
			ctx = WithTracer(ctx, msg.tracer)

			pipeline, err = BuildFlow(ctx, w.configStorage, msg.pipelineName, &msg.tracer.Flow)
			if err != nil {
				msg.replyTo <- err
				return
			}

			if err := pipeline.Init(ctx); err != nil {
				_ = pipeline.Close(ctx)
				msg.replyTo <- err
				return
			}
		} else {
			pipeline, err = w.getOrBuildPipeline(ctx, msg.pipelineName)
		}

		if err != nil {
			msg.replyTo <- err
			return
		}

		err = pipeline.Process(ctx, msg.entrypoint, msg.record)

		if msg.tracer != nil {
			_ = pipeline.Close(ctx)
		}

		msg.replyTo <- err
	}()
}

func (msg invalidateCacheMessage) handle(ctx context.Context, w *worker) {
	w.cacheMu.Lock()
	defer w.cacheMu.Unlock()

	pipeline, ok := w.cache[msg.pipelineName]
	if ok {
		msg.replyTo <- pipeline.Close(ctx)
		delete(w.cache, msg.pipelineName)
	} else {
		msg.replyTo <- nil
	}
}

func (msg invalidateAllCacheMessage) handle(ctx context.Context, w *worker) {
	w.cacheMu.Lock()
	defer w.cacheMu.Unlock()

	var errs []error

	for name, pipeline := range w.cache {
		if err := pipeline.Close(ctx); err != nil {
			errs = append(errs, err)
		}
		delete(w.cache, name)
	}

	if len(errs) > 0 {
		msg.replyTo <- errors.Join(errs...)
	} else {
		msg.replyTo <- nil
	}
}
