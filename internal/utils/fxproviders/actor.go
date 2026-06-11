package fxproviders

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"
)

func ProvideActor[T actor.Actor](constructor any) fx.Option {
	type in struct {
		fx.In
		A T
	}

	return fx.Options(
		fx.Provide(constructor),
		fx.Invoke(func(lc fx.Lifecycle, p in) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					p.A.Start()
					return nil
				},
				OnStop: func(context.Context) error {
					p.A.Stop()
					return nil
				},
			})
		}),
	)
}
