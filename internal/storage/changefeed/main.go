package changefeed

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/utils/fxproviders"
)

type Notifier interface {
	actor.Actor

	Subscribe(ctx context.Context) (actor.MailboxReceiver[ChangeEvent], error)
	Notify(ctx context.Context, event ChangeEvent) error
}

type notifierImpl struct {
	actor.Actor

	subMbox   actor.MailboxSender[subscribeMessage]
	eventMbox actor.MailboxSender[ChangeEvent]
}

var _ Notifier = (*notifierImpl)(nil)

func NewNotifier() fx.Option {
	return fx.Module(
		"storage.changefeed",
		fxproviders.ProvideMailbox[subscribeMessage](),
		fxproviders.ProvideMailbox[ChangeEvent](),
		fxproviders.ProvideActor[Notifier](func(
			subMbox actor.Mailbox[subscribeMessage],
			eventMbox actor.Mailbox[ChangeEvent],
		) Notifier {
			return &notifierImpl{
				Actor: actor.New(&worker{
					subscribers: make(subscriberSet),
					subMbox:     subMbox,
					eventMbox:   eventMbox,
				}),

				subMbox:   subMbox,
				eventMbox: eventMbox,
			}
		}),
	)
}

func (n *notifierImpl) Subscribe(ctx context.Context) (actor.MailboxReceiver[ChangeEvent], error) {
	eventM := actor.NewMailbox[ChangeEvent]()
	eventM.Start()

	readyC := make(chan error, 1)
	msg := subscribeMessage{
		senderM: eventM,
		readyC:  readyC,
		doneC:   ctx.Done(),
	}

	if err := n.subMbox.Send(ctx, msg); err != nil {
		eventM.Stop()
		return nil, err
	}

	if err := <-readyC; err != nil {
		eventM.Stop()
		return nil, err
	}

	return eventM, nil
}

func (n *notifierImpl) Notify(ctx context.Context, event ChangeEvent) error {
	return n.eventMbox.Send(ctx, event)
}
