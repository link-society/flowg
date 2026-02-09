package lognotify

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/utils/fxproviders"

	"link-society.com/flowg/internal/models"
)

type LogNotifier interface {
	actor.Actor

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
		fxproviders.ProvideMailbox[SubscribeMessage](),
		fxproviders.ProvideMailbox[LogMessage](),
		fxproviders.ProvideActor[LogNotifier](func(
			subMbox actor.Mailbox[SubscribeMessage],
			logMbox actor.Mailbox[LogMessage],
		) LogNotifier {
			return &logNotifierImpl{
				Actor: actor.New(&worker{
					subscribers: make(map[string]subscriberSet),
					subMbox:     subMbox,
					logMbox:     logMbox,
				}),

				subMbox: subMbox,
				logMbox: logMbox,
			}
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
