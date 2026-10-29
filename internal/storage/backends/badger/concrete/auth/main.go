package auth

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"go.link-society.com/flowg/internal/storage/backends/badger"
	"go.link-society.com/flowg/internal/storage/databases/auth"
	storage "go.link-society.com/flowg/internal/storage/interfaces"
)

// Options configures the badger-backed authentication storage.
type Options struct {
	Directory string
	InMemory  bool
	ReadOnly  bool
}

type deps struct {
	fx.In

	Adapter *badger.BadgerAdapter `name:"storage.auth"`
}

// DefaultOptions returns the default [Options] for the authentication storage.
func DefaultOptions() Options {
	return Options{
		Directory: "",
		InMemory:  false,
		ReadOnly:  false,
	}
}

// adapterOptions returns the options of the underlying badger adapter.
func (o Options) adapterOptions() badger.AdapterOptions {
	return badger.AdapterOptions{
		Directory: o.Directory,
		InMemory:  o.InMemory,
		ReadOnly:  o.ReadOnly,
	}
}

// Module returns an fx module providing a badger-backed [storage.AuthStorage]
// configured by the [Options] provided to the graph.
func Module() fx.Option {
	return fx.Module(
		"storage.auth",
		badger.NewAdapter("storage.auth", Options.adapterOptions),
		fx.Provide(func(lc fx.Lifecycle, d deps) storage.AuthStorage {
			storage := auth.NewStorage(d.Adapter)

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := migrateAlertScopes(ctx, d.Adapter); err != nil {
						return fmt.Errorf("failed to migrate alerts: %w", err)
					}

					return nil
				},
			})

			return storage
		}),
	)
}
