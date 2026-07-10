package config

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"link-society.com/flowg/internal/storage/backends/badger"
	"link-society.com/flowg/internal/storage/databases/config"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// Options configures the badger-backed configuration storage.
type Options struct {
	Directory string
	InMemory  bool
	ReadOnly  bool
}

type deps struct {
	fx.In

	Adapter *badger.BadgerAdapter `name:"storage.config"`
}

// DefaultOptions returns the default [Options] for the configuration storage.
func DefaultOptions() Options {
	return Options{
		Directory: "",
		InMemory:  false,
		ReadOnly:  false,
	}
}

// NewStorage returns an fx module providing a badger-backed
// [storage.ConfigStorage] configured with the given options.
func NewStorage(opts Options) fx.Option {
	adapterOpts := badger.AdapterOptions{
		LogChannel: "storage.config",
		Directory:  opts.Directory,
		InMemory:   opts.InMemory,
		ReadOnly:   opts.ReadOnly,
	}

	return fx.Module(
		"storage.config",
		badger.NewAdapter(adapterOpts),
		fx.Provide(func(lc fx.Lifecycle, d deps) storage.ConfigStorage {
			storage := config.NewStorage(d.Adapter)

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := migrateAlerts(opts.Directory); err != nil {
						return fmt.Errorf("failed to migrate alerts: %w", err)
					}

					if err := migrateToBadger(ctx, opts.Directory, storage); err != nil {
						return fmt.Errorf("failed to migrate to badger: %w", err)
					}

					return nil
				},
			})

			return storage
		}),
	)
}
