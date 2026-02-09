package fxproviders

import (
	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"
)

func ProvideMailbox[T any]() fx.Option {
	return ProvideActor[actor.Mailbox[T]](func() actor.Mailbox[T] {
		return actor.NewMailbox[T]()
	})
}
