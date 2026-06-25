package pipelines

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage"

	"link-society.com/flowg/internal/engines/lognotify"
)

// Runner executes pipelines against incoming log records. It is the public entry
// point of the engine: callers submit a record to a named pipeline and the
// runner drives it through the pipeline's node graph.
type Runner interface {
	// Run pushes a record through pipelineName, starting at entrypoint (e.g.
	// "direct" or "syslog"), and blocks until processing completes.
	Run(ctx context.Context, pipelineName string, entrypoint string, record *models.LogRecord) error
	// InvalidateCachedBuild drops the compiled build of a single pipeline so the
	// next Run rebuilds it from storage; call it after the pipeline changes.
	InvalidateCachedBuild(ctx context.Context, pipelineName string) error
	// InvalidateAllCachedBuilds drops every compiled pipeline, e.g. on shutdown
	// or after a bulk configuration import.
	InvalidateAllCachedBuilds(ctx context.Context) error
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
			a := actor.New(&worker{
				mbox:          mbox,
				configStorage: d.ConfigStorage,
				logStorage:    d.LogStorage,
				logNotifier:   d.LogNotifier,
				cache:         make(map[string]*Pipeline),
			})
			runner := &runnerImpl{mbox: mbox}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					a.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					if err := runner.InvalidateAllCachedBuilds(ctx); err != nil {
						return err
					}

					a.Stop()
					return nil
				},
			})

			return runner
		}),
	)
}

// Run sends a processing request to the actor and waits for the result. The
// active tracer (if any, set via WithTracer) is forwarded so dry runs can record
// per-node traces.
func (r *runnerImpl) Run(
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

// InvalidateCachedBuild asks the actor to evict and close the cached build of a
// single pipeline.
func (r *runnerImpl) InvalidateCachedBuild(
	ctx context.Context,
	pipelineName string,
) error {
	replyTo := make(chan error)

	err := r.mbox.Send(ctx, invalidateCacheMessage{
		replyTo:      replyTo,
		pipelineName: pipelineName,
	})
	if err != nil {
		return err
	}

	return <-replyTo
}

// InvalidateAllCachedBuilds asks the actor to evict and close every cached
// pipeline build.
func (r *runnerImpl) InvalidateAllCachedBuilds(
	ctx context.Context,
) error {
	replyTo := make(chan error)

	err := r.mbox.Send(ctx, invalidateAllCacheMessage{
		replyTo: replyTo,
	})
	if err != nil {
		return err
	}

	return <-replyTo
}
