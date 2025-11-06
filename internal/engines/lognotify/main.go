package lognotify

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/models"
)

type LogNotifier interface {
	Subscribe(ctx context.Context, stream string) (actor.MailboxReceiver[LogMessage], error)
	Notify(ctx context.Context, stream string, logKey string, logRecord models.LogRecord) error
}

type logNotifierImpl struct {
	actor.Actor

	subMbox actor.MailboxSender[SubscribeMessage]
	logMbox actor.MailboxSender[LogMessage]
}

var _ LogNotifier = (*logNotifierImpl)(nil)

func NewLogNotifier() fx.Option {
	return fx.Module(
		"lognotifier",
		fx.Provide(func(lc fx.Lifecycle) actor.Mailbox[SubscribeMessage] {
			mbox := actor.NewMailbox[SubscribeMessage]()

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
		fx.Provide(func(lc fx.Lifecycle) actor.Mailbox[LogMessage] {
			mbox := actor.NewMailbox[LogMessage]()

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
		fx.Provide(func(
			lc fx.Lifecycle,
			subMbox actor.Mailbox[SubscribeMessage],
			logMbox actor.Mailbox[LogMessage],
		) LogNotifier {
			logNotifier := &logNotifierImpl{
				Actor: actor.New(&worker{
					subscribers: make(map[string]subscriberSet),
					subMbox:     subMbox,
					logMbox:     logMbox,
				}),

				subMbox: subMbox,
				logMbox: logMbox,
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logNotifier.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logNotifier.Stop()
					return nil
				},
			})

			return logNotifier
		}),
	)
}

func (n *logNotifierImpl) Subscribe(ctx context.Context, stream string) (actor.MailboxReceiver[LogMessage], error) {
	logM := actor.NewMailbox[LogMessage]()
	logM.Start()

	readyC := make(chan ReadyResponse, 1)
	msg := SubscribeMessage{
		Stream:  stream,
		SenderM: logM,
		ReadyC:  readyC,
		DoneC:   ctx.Done(),
	}

	err := n.subMbox.Send(ctx, msg)
	if err != nil {
		logM.Stop()
		return nil, err
	}

	resp := <-readyC
	if resp.Err != nil {
		logM.Stop()
		return nil, resp.Err
	}

	return logM, nil
}

func (n *logNotifierImpl) Notify(ctx context.Context, stream string, logKey string, logRecord models.LogRecord) error {
	msg := LogMessage{
		Stream:    stream,
		LogKey:    logKey,
		LogRecord: logRecord,
	}
	return n.logMbox.Send(ctx, msg)
}
