package config

import (
	"context"
	"log/slog"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"go.link-society.com/flowg/internal/engines/confignotify"

	"go.link-society.com/flowg/internal/storage/backends/foundation"
	"go.link-society.com/flowg/internal/storage/databases/config"
	storage "go.link-society.com/flowg/internal/storage/interfaces"
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

	Logger  *slog.Logger
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

// adapterOptions returns the options of the underlying FoundationDB adapter.
func (o Options) adapterOptions() foundation.AdapterOptions {
	return foundation.AdapterOptions{
		ClusterFile:      o.ClusterFile,
		ConnectionString: o.ConnectionString,
		KeySpace:         o.KeySpace,
		Namespace:        "config",
	}
}

// Module returns an fx module providing a FoundationDB-backed
// [storage.ConfigStorage] configured by the [Options] provided to the graph.
//
// The module also starts a changelog watcher that translates config mutations
// (local or from other nodes) into confignotify events.
func Module() fx.Option {
	return fx.Module(
		"storage.config",
		foundation.NewAdapter("storage.config", Options.adapterOptions),
		fx.Provide(func(d deps) *config.Storage[*foundation.FoundationQueryTx, *foundation.FoundationMutationTx] {
			// no local notifier: change events come from the key watcher, so every
			// node (including the writer) observes mutations through the same path
			return config.NewStorage(d.Logger, d.Adapter, nil)
		}),
		fx.Provide(func(s *config.Storage[*foundation.FoundationQueryTx, *foundation.FoundationMutationTx]) storage.ConfigStorage {
			return s
		}),
		fx.Invoke(func(
			lc fx.Lifecycle,
			d deps,
			s *config.Storage[*foundation.FoundationQueryTx, *foundation.FoundationMutationTx],
		) {
			if d.Notifier == nil {
				return
			}

			watcher := actor.New(&watchWorker{
				logger:      d.Logger.With(slog.String("channel", "storage.config")),
				adapter:     d.Adapter,
				notifier:    d.Notifier,
				systemCache: s,
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
