package confignotify

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"
)

// Notifier is the fan-out bus for configuration changes. The config storage
// publishes an event after every successful mutation, and interested engines
// (the pipeline runner) subscribe to react to them. On clustered backends the
// storage watcher publishes the events instead, so subscribers observe local
// and remote changes through the same bus.
type Notifier interface {
	// Subscribe registers the caller as a listener and returns a mailbox that
	// receives every subsequent event. The subscription is torn down when ctx is
	// cancelled.
	Subscribe(ctx context.Context) (actor.MailboxReceiver[Event], error)
	// Notify broadcasts a configuration change to every current subscriber; it
	// is a no-op when nobody is listening.
	Notify(ctx context.Context, event Event) error
}

type notifierImpl struct {
	actor.Actor

	subMbox   actor.MailboxSender[SubscribeMessage]
	eventMbox actor.MailboxSender[Event]
}

var _ Notifier = (*notifierImpl)(nil)

// NewNotifier returns an fx module providing a Notifier backed by a single
// actor. The actor and its two mailboxes (subscriptions and event broadcasts)
// are started and stopped with the application lifecycle.
func NewNotifier() fx.Option {
	return fx.Module(
		"confignotifier",
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
		fx.Provide(func(lc fx.Lifecycle) actor.Mailbox[Event] {
			mbox := actor.NewMailbox[Event]()

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
			eventMbox actor.Mailbox[Event],
		) Notifier {
			notifier := &notifierImpl{
				Actor: actor.New(&worker{
					subscribers: make(map[actor.Mailbox[Event]]<-chan struct{}),
					subMbox:     subMbox,
					eventMbox:   eventMbox,
				}),

				subMbox:   subMbox,
				eventMbox: eventMbox,
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					notifier.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					notifier.Stop()
					return nil
				},
			})

			return notifier
		}),
	)
}

// Subscribe asks the actor to register a new per-subscriber mailbox and blocks
// until the actor confirms the registration, guaranteeing no event is missed
// between the call and the first delivery.
func (n *notifierImpl) Subscribe(ctx context.Context) (actor.MailboxReceiver[Event], error) {
	eventM := actor.NewMailbox[Event]()
	eventM.Start()

	readyC := make(chan ReadyResponse, 1)
	msg := SubscribeMessage{
		SenderM: eventM,
		ReadyC:  readyC,
		DoneC:   ctx.Done(),
	}

	err := n.subMbox.Send(ctx, msg)
	if err != nil {
		eventM.Stop()
		return nil, err
	}

	resp := <-readyC
	if resp.Err != nil {
		eventM.Stop()
		return nil, resp.Err
	}

	return eventM, nil
}

// Notify hands an event to the actor for broadcasting; it returns as soon as
// the event is queued, without waiting for delivery to subscribers.
func (n *notifierImpl) Notify(ctx context.Context, event Event) error {
	return n.eventMbox.Send(ctx, event)
}
