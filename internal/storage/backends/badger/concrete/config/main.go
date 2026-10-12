package config

import (
	"context"
	"fmt"
	"log/slog"

	"go.uber.org/fx"

	"link-society.com/flowg/internal/engines/confignotify"

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

	Logger  *slog.Logger
	Options Options
	Adapter *badger.BadgerAdapter `name:"storage.config"`

	// Notifier broadcasts config mutations to the pipeline runner; it is
	// optional so the storage can run standalone (e.g. in tests).
	Notifier confignotify.Notifier `optional:"true"`
}

// DefaultOptions returns the default [Options] for the configuration storage.
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

// Module returns an fx module providing a badger-backed [storage.ConfigStorage]
// configured by the [Options] provided to the graph.
func Module() fx.Option {
	return fx.Module(
		"storage.config",
		badger.NewAdapter("storage.config", Options.adapterOptions),
		fx.Provide(func(lc fx.Lifecycle, d deps) storage.ConfigStorage {
			storage := config.NewStorage(d.Logger, d.Adapter, d.Notifier)

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := migrateAlerts(d.Options.Directory); err != nil {
						return fmt.Errorf("failed to migrate alerts: %w", err)
					}

					if err := migrateToBadger(ctx, d.Options.Directory, storage); err != nil {
						return fmt.Errorf("failed to migrate to badger: %w", err)
					}

					return nil
				},
			})

			return storage
		}),
	)
}
