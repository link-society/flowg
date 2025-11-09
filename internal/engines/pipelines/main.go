package pipelines

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
)

type Runner interface {
	Run(ctx context.Context, pipelineName string, entrypoint string, record *models.LogRecord) error
	InvalidateCachedBuild(ctx context.Context, pipelineName string) error
	InvalidateAllCachedBuilds(ctx context.Context) error
}

type runnerImpl struct {
	mbox actor.MailboxSender[message]
}

type deps struct {
	fx.In

	ConfigStorage config.Storage
	LogStorage    log.Storage
	LogNotifier   lognotify.LogNotifier
}

var _ Runner = (*runnerImpl)(nil)

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
	})
	if err != nil {
		return err
	}

	return <-replyTo
}

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
