package config

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/engines/confignotify"

	"link-society.com/flowg/internal/storage/backends/foundation"
	"link-society.com/flowg/internal/storage/databases/config"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// Options configures the FoundationDB-backed configuration storage.
type Options struct {
	// ClusterFile is the path to the fdb.cluster file; empty uses the default.
	ClusterFile string
	// ConnectionString is the FoundationDB connection string; empty uses the default.
	ConnectionString string
	// KeySpace is the root prefix shared by every FlowG storage.
	KeySpace string
}

type deps struct {
	fx.In

	Adapter *foundation.FoundationAdapter `name:"storage.config"`

	// Notifier receives the events emitted by the changelog watcher; it is
	// optional so the storage can run standalone (e.g. in tests).
	Notifier confignotify.Notifier `optional:"true"`
}

// DefaultOptions returns the default [Options] for the configuration storage.
func DefaultOptions() Options {
	return Options{
		ClusterFile:      "",
		ConnectionString: "",
		KeySpace:         "flowg",
	}
}

// NewStorage returns an fx module providing a FoundationDB-backed
// [storage.ConfigStorage] configured with the given options.
//
// The module also starts a changelog watcher that translates config mutations
// (local or from other nodes) into confignotify events.
func NewStorage(opts Options) fx.Option {
	adapterOpts := foundation.AdapterOptions{
		LogChannel:       "storage.config",
		ClusterFile:      opts.ClusterFile,
		ConnectionString: opts.ConnectionString,
		KeySpace:         opts.KeySpace,
		Namespace:        "config",
		EnableChangeLog:  true,
	}

	return fx.Module(
		"storage.config",
		foundation.NewAdapter(adapterOpts),
		fx.Provide(func(d deps) storage.ConfigStorage {
			// no local notifier: change events come from the key watcher, so every
			// node (including the writer) observes mutations through the same path
			return config.NewStorage(d.Adapter, nil)
		}),
		fx.Invoke(func(lc fx.Lifecycle, d deps) {
			if d.Notifier == nil {
				return
			}

			watcher := actor.New(&watchWorker{
				adapter:  d.Adapter,
				notifier: d.Notifier,
			})

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					watcher.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					watcher.Stop()
					return nil
				},
			})
		}),
	)
}
