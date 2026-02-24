package pipelines

import (
	"context"
	"errors"

	"link-society.com/flowg/internal/models"
)

const (
	DIRECT_ENTRYPOINT = "direct"
	SYSLOG_ENTRYPOINT = "syslog"
)

type message interface {
	handle(ctx context.Context, w *worker)
}

type logMessage struct {
	replyTo chan<- error

	pipelineName string
	entrypoint   string
	record       *models.LogRecord
	tracer       *NodeTracer
}

type invalidateCacheMessage struct {
	replyTo      chan<- error
	pipelineName string
}

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
		} else {
			pipeline, err = w.getOrBuildPipeline(ctx, msg.pipelineName)
		}

		if err != nil {
			msg.replyTo <- err
			return
		}

		err = pipeline.Process(ctx, msg.entrypoint, msg.record)
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
