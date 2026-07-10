package log

import (
	"context"

	"time"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/storage/backends/badger"
	"link-society.com/flowg/internal/storage/databases/log"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// Options configures the badger-backed log storage.
type Options struct {
	Directory  string
	InMemory   bool
	ReadOnly   bool
	GCInterval time.Duration
}

type deps struct {
	fx.In

	Adapter *badger.BadgerAdapter `name:"storage.log"`
}

// DefaultOptions returns the default [Options] for the log storage.
func DefaultOptions() Options {
	return Options{
		Directory:  "",
		InMemory:   false,
		ReadOnly:   false,
		GCInterval: 5 * time.Minute,
	}
}

// NewStorage returns an fx module providing a badger-backed [storage.LogStorage]
// configured with the given options. It also starts a background worker that
// periodically runs the value-log garbage collector.
func NewStorage(opts Options) fx.Option {
	adapterOpts := badger.AdapterOptions{
		LogChannel: "storage.log",
		Directory:  opts.Directory,
		InMemory:   opts.InMemory,
		ReadOnly:   opts.ReadOnly,
	}

	return fx.Module(
		"storage.log",
		badger.NewAdapter(adapterOpts),
		fx.Provide(func(lc fx.Lifecycle, d deps) storage.LogStorage {
			storage := log.NewStorage(d.Adapter)

			gc := actor.New(log.NewGarbageCollector(d.Adapter, opts.GCInterval))

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					gc.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					gc.Stop()
					return nil
				},
			})

			return storage
		}),
	)
}
