package auth

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"link-society.com/flowg/internal/storage/backends/badger"
	"link-society.com/flowg/internal/storage/databases/auth"
	storage "link-society.com/flowg/internal/storage/interfaces"
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

// NewStorage returns an fx module providing a badger-backed
// [storage.AuthStorage] configured with the given options.
func NewStorage(opts Options) fx.Option {
	adapterOpts := badger.AdapterOptions{
		LogChannel: "storage.auth",
		Directory:  opts.Directory,
		InMemory:   opts.InMemory,
		ReadOnly:   opts.ReadOnly,
	}

	return fx.Module(
		"storage.auth",
		badger.NewAdapter(adapterOpts),
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
