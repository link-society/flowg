package pipelines

import (
	"context"

	"io"

	"github.com/prometheus/common/expfmt"

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

// scrapMetricsMessage forwards a request to scrap metrics from a pipeline.
type scrapMetricsMessage struct {
	replyTo chan<- error

	pipelineName string
	w            io.Writer
}

// runMessage requests that a pipeline be loaded from storage, compiled and
// initialized.
type runMessage struct {
	replyTo chan<- error

	pipelineName string
}

var (
	_ message = logMessage{}
	_ message = scrapMetricsMessage{}
	_ message = runMessage{}
	_ message = terminateMessage{}
)

// terminateMessage requests that a pipeline be terminated and removed.
type terminateMessage struct {
	replyTo chan<- error

	pipelineName string
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
			pipeline, err = w.getPipeline(ctx, msg.pipelineName)
			if err != nil {
				msg.replyTo <- err
				return
			}
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

func (msg scrapMetricsMessage) handle(ctx context.Context, w *worker) {
	go func() {
		var pipeline *Pipeline
		var err error
		defer close(msg.replyTo)

		pipeline, err = w.getPipeline(ctx, msg.pipelineName)
		if err != nil {
			msg.replyTo <- err
			return
		}

		families, err := pipeline.Metrics.Gather()
		if err != nil {
			msg.replyTo <- err
			return
		}

		encoder := expfmt.NewEncoder(msg.w, expfmt.NewFormat(expfmt.TypeTextPlain))

		for _, family := range families {
			if err := encoder.Encode(family); err != nil {
				msg.replyTo <- err
				return
			}
		}

		msg.replyTo <- nil
	}()
}

func (msg runMessage) handle(ctx context.Context, w *worker) {
	defer close(msg.replyTo)
	msg.replyTo <- w.buildPipeline(ctx, msg.pipelineName)
}

func (msg terminateMessage) handle(ctx context.Context, w *worker) {
	defer close(msg.replyTo)
	msg.replyTo <- w.closePipeline(ctx, msg.pipelineName)
}
