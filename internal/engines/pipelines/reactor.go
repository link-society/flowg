package pipelines

import (
	"context"
	"log/slog"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/utils/fxproviders"

	"link-society.com/flowg/internal/storage/changefeed"
	"link-society.com/flowg/internal/storage/config"
)

type reactor struct {
	actor.Actor
}

type reactorDeps struct {
	fx.In

	Runner   Runner
	Notifier changefeed.Notifier
}

func NewReactor() fx.Option {
	return fx.Module(
		"pipelineReactor",
		fxproviders.ProvideActor[*reactor](func(d reactorDeps) *reactor {
			return &reactor{
				Actor: actor.New(&reactorWorker{
					runner:   d.Runner,
					notifier: d.Notifier,
				}),
			}
		}),
	)
}

type reactorWorker struct {
	runner   Runner
	notifier changefeed.Notifier

	eventR actor.MailboxReceiver[changefeed.ChangeEvent]
}

var _ actor.Worker = (*reactorWorker)(nil)

func (w *reactorWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	if w.eventR == nil {
		eventR, err := w.notifier.Subscribe(ctx)
		if err != nil {
			if ctx.Err() == nil {
				slog.ErrorContext(
					ctx,
					"failed to subscribe to change feed",
					slog.String("channel", "pipelines.reactor"),
					slog.String("error", err.Error()),
				)
			}

			return actor.WorkerEnd
		}

		w.eventR = eventR
	}

	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case event, ok := <-w.eventR.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		w.handle(ctx, event)
		return actor.WorkerContinue
	}
}

func (w *reactorWorker) handle(ctx context.Context, event changefeed.ChangeEvent) {
	if event.Namespace != changefeed.NamespaceConfig {
		return
	}

	if !event.Resync && event.Kind == config.SystemItemType {
		return
	}

	if !event.Resync && event.Kind == config.PipelineItemType {
		if err := w.runner.InvalidateCachedBuild(ctx, event.Name); err != nil {
			slog.ErrorContext(
				ctx,
				"failed to invalidate pipeline cache",
				slog.String("channel", "pipelines.reactor"),
				slog.String("pipeline", event.Name),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if err := w.runner.InvalidateAllCachedBuilds(ctx); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to invalidate pipeline cache",
			slog.String("channel", "pipelines.reactor"),
			slog.String("error", err.Error()),
		)
	}
}
