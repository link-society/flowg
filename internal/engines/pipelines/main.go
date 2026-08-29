package pipelines

import (
	"context"
	"errors"

	"fmt"
	"log/slog"

	"io"
	"sync"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/models"

	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/engines/confignotify"
	"link-society.com/flowg/internal/engines/lognotify"
)

// Runner executes pipelines against incoming log records. It is the public entry
// point of the engine: callers submit a record to a named pipeline and the
// runner drives it through the pipeline's node graph.
//
// The pipeline lifecycle itself is not part of the interface: builds are
// started eagerly at boot and follow configuration changes broadcast on the
// confignotify bus.
type Runner interface {
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

	ConfigStorage  storage.ConfigStorage
	LogStorage     storage.LogStorage
	LogNotifier    lognotify.LogNotifier
	ConfigNotifier confignotify.Notifier
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
					// the subscription must outlive the start phase; it is torn down
					// with the notifier's own lifecycle
					events, err := d.ConfigNotifier.Subscribe(context.Background())
					if err != nil {
						return fmt.Errorf("failed to subscribe to config changes: %w", err)
					}
					w.eventsC = events.ReceiveC()

					a.Start()

					pipelines, err := d.ConfigStorage.ListPipelines(ctx)
					if err != nil {
						return fmt.Errorf("failed to start pipelines: %w", err)
					}

					// a broken pipeline must not prevent the server (and the UI needed
					// to fix it) from starting
					for _, pipelineName := range pipelines {
						if err := runner.run(ctx, pipelineName); err != nil {
							slog.ErrorContext(
								ctx,
								"failed to start pipeline",
								"channel", "pipelines",
								"pipeline", pipelineName,
								"error", err.Error(),
							)
						}
					}

					return nil
				},
				OnStop: func(ctx context.Context) error {
					w.cacheMu.Lock()
					names := make([]string, 0, len(w.cache))
					for pipelineName := range w.cache {
						names = append(names, pipelineName)
					}
					w.cacheMu.Unlock()

					var errs []error
					for _, pipelineName := range names {
						if err := runner.terminate(ctx, pipelineName); err != nil {
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

// run sends a request to the actor to build and start a pipeline. Lifecycle is
// internal: it is driven by the fx hooks and the confignotify events.
func (r *runnerImpl) run(ctx context.Context, pipelineName string) error {
	replyTo := make(chan error, 1)

	err := r.mbox.Send(ctx, runMessage{
		replyTo:      replyTo,
		pipelineName: pipelineName,
	})
	if err != nil {
		return err
	}

	select {
	case err := <-replyTo:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// terminate sends a request to the actor to terminate a pipeline. It waits for
// in-flight records to drain, bounded by ctx; the teardown itself always
// completes in the background.
func (r *runnerImpl) terminate(ctx context.Context, pipelineName string) error {
	replyTo := make(chan error, 1)

	err := r.mbox.Send(ctx, terminateMessage{
		replyTo:      replyTo,
		pipelineName: pipelineName,
	})
	if err != nil {
		return err
	}

	select {
	case err := <-replyTo:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
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
	replyTo := make(chan error, 1)

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

	select {
	case err := <-replyTo:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ScrapMetrics sends a request to the actor to scrap metrics from a pipeline.
func (r *runnerImpl) ScrapMetrics(
	ctx context.Context,
	pipelineName string,
	w io.Writer,
) error {
	replyTo := make(chan error, 1)

	err := r.mbox.Send(ctx, scrapMetricsMessage{
		replyTo:      replyTo,
		pipelineName: pipelineName,
		w:            w,
	})
	if err != nil {
		return err
	}

	select {
	case err := <-replyTo:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
