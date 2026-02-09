package pipelines

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/utils/fxproviders"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
)

type Runner interface {
	actor.Actor

	Run(ctx context.Context, pipelineName string, entrypoint string, record *models.LogRecord) error
	InvalidateCachedBuild(ctx context.Context, pipelineName string) error
	InvalidateAllCachedBuilds(ctx context.Context) error
}

type runnerImpl struct {
	actor.Actor

	mbox actor.MailboxSender[message]
}

type deps struct {
	fx.In

	Mailbox actor.Mailbox[message]

	ConfigStorage config.Storage
	LogStorage    log.Storage
	LogNotifier   lognotify.LogNotifier
}

var _ Runner = (*runnerImpl)(nil)

func NewRunner() fx.Option {
	return fx.Module(
		"pipelineRunner",
		fxproviders.ProvideMailbox[message](),
		fxproviders.ProvideActor[Runner](
			func(d deps) Runner {
				w := &worker{
					mbox:          d.Mailbox,
					configStorage: d.ConfigStorage,
					logStorage:    d.LogStorage,
					logNotifier:   d.LogNotifier,
					cache:         make(map[string]*Pipeline),
				}
				return &runnerImpl{
					Actor: actor.New(w),
					mbox:  d.Mailbox,
				}
			},
		),
		fx.Invoke(func(lc fx.Lifecycle, runner Runner) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return runner.InvalidateAllCachedBuilds(ctx)
				},
			})
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
		tracer:       GetTracer(ctx),
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
