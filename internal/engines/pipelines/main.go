package pipelines

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"io"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/models"

	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/engines/lognotify"
)

// Runner executes pipelines against incoming log records. It is the public entry
// point of the engine: callers submit a record to a named pipeline and the
// runner drives it through the pipeline's node graph.
type Runner interface {
	// Run compiles and starts a pipeline
	Run(ctx context.Context, pipelineName string) error
	// Terminate stops a running pipeline.
	Terminate(ctx context.Context, pipelineName string) error
	// Process pushes a record through pipelineName, starting at entrypoint (e.g.
	// "direct" or "syslog"), and blocks until processing completes.
	Process(ctx context.Context, pipelineName string, entrypoint string, record *models.LogRecord) error
	// ScrapMetrics forwards a request to scrap metrics from a pipeline.
	ScrapMetrics(ctx context.Context, pipelineName string, w io.Writer) error
}

type runnerImpl struct {
	mbox actor.MailboxSender[message]
}

type deps struct {
	fx.In

	ConfigStorage storage.ConfigStorage
	LogStorage    storage.LogStorage
	LogNotifier   lognotify.LogNotifier
}

var _ Runner = (*runnerImpl)(nil)

// NewRunner returns an fx module providing a Runner backed by a single actor.
// The actor owns the pipeline build cache; on shutdown every cached build is
// invalidated (closing forwarders and transformers) before the actor stops.
func NewRunner() fx.Option {
	return fx.Module(
		"pipelineRunner",
		fx.Provide(func(lc fx.Lifecycle) actor.Mailbox[message] {
			mbox := actor.NewMailbox[message]()

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					mbox.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					mbox.Stop()
					return nil
				},
			})

			return mbox
		}),
		fx.Provide(func(lc fx.Lifecycle, d deps, mbox actor.Mailbox[message]) Runner {
			w := &worker{
				mbox:          mbox,
				configStorage: d.ConfigStorage,
				logStorage:    d.LogStorage,
				logNotifier:   d.LogNotifier,

				cache:   make(map[string]*Pipeline),
				cacheMu: sync.Mutex{},
			}
			a := actor.New(w)
			runner := &runnerImpl{mbox: mbox}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					a.Start()

					pipelines, err := d.ConfigStorage.ListPipelines(ctx)
					if err != nil {
						return fmt.Errorf("failed to start pipelines: %w", err)
					}

					var errs []error
					for _, pipelineName := range pipelines {
						if err := runner.Run(ctx, pipelineName); err != nil {
							errs = append(errs, err)
						}
					}
					if len(errs) > 0 {
						return fmt.Errorf("failed to start pipelines: %w", errors.Join(errs...))
					}

					return nil
				},
				OnStop: func(ctx context.Context) error {
					var errs []error
					for pipelineName := range w.cache {
						if err := runner.Terminate(ctx, pipelineName); err != nil {
							errs = append(errs, err)
						}
					}
					if len(errs) > 0 {
						return fmt.Errorf("failed to close pipelines: %w", errors.Join(errs...))
					}

					a.Stop()

					return nil
				},
			})

			return runner
		}),
	)
}

// Run sends a request to the actor to run a pipeline.
func (r *runnerImpl) Run(ctx context.Context, pipelineName string) error {
	replyTo := make(chan error)

	err := r.mbox.Send(ctx, runMessage{
		replyTo:      replyTo,
		pipelineName: pipelineName,
	})
	if err != nil {
		return err
	}

	return <-replyTo
}

// Terminate sends a request to the actor to terminate a pipeline.
func (r *runnerImpl) Terminate(ctx context.Context, pipelineName string) error {
	replyTo := make(chan error)

	err := r.mbox.Send(ctx, terminateMessage{
		replyTo:      replyTo,
		pipelineName: pipelineName,
	})
	if err != nil {
		return err
	}

	return <-replyTo
}

// Process sends a processing request to the actor and waits for the result. The
// active tracer (if any, set via WithTracer) is forwarded so dry runs can record
// per-node traces.
func (r *runnerImpl) Process(
	ctx context.Context,
	pipelineName string,
	entrypoint string,
	record *models.LogRecord,
) error {
	replyTo := make(chan error)

	err := r.mbox.Send(ctx, logMessage{
		replyTo: replyTo,

		pipelineName: pipelineName,
		entrypoint:   entrypoint,
		record:       record,
		tracer:       GetTracer(ctx),
	})
	if err != nil {
		return err
	}

	return <-replyTo
}

// ScrapMetrics sends a request to the actor to scrap metrics from a pipeline.
func (r *runnerImpl) ScrapMetrics(
	ctx context.Context,
	pipelineName string,
	w io.Writer,
) error {
	replyTo := make(chan error)

	err := r.mbox.Send(ctx, scrapMetricsMessage{
		replyTo:      replyTo,
		pipelineName: pipelineName,
		w:            w,
	})
	if err != nil {
		return err
	}

	return <-replyTo
}
