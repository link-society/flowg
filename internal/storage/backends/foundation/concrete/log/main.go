package log

import (
	"context"
	"log/slog"
	"time"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"go.link-society.com/flowg/internal/models"

	"go.link-society.com/flowg/internal/storage/backends/foundation"

	"go.link-society.com/flowg/internal/storage/databases/log"
	storage "go.link-society.com/flowg/internal/storage/interfaces"
)

// Options configures the FoundationDB-backed log storage.
type Options struct {
	// ClusterFile is the path to the fdb.cluster file; empty uses the default.
	ClusterFile string
	// ConnectionString is the FoundationDB connection string; empty uses the default.
	ConnectionString string
	// KeySpace is the root prefix shared by every FlowG storage.
	KeySpace string
	// GCInterval is how often the retention and TTL garbage collectors run.
	GCInterval time.Duration
	// BatchSize bounds how many keys whole-stream operations touch per
	// transaction; zero uses the storage default.
	BatchSize int
	// StreamConfigOverride, when not nil, is used as the configuration of every
	// stream instead of the stored one.
	StreamConfigOverride *models.StreamConfig
}

type deps struct {
	fx.In

	Logger  *slog.Logger
	Options Options
	Adapter *foundation.FoundationAdapter `name:"storage.log"`
}

// DefaultOptions returns the default [Options] for the log storage.
func DefaultOptions() Options {
	return Options{
		ClusterFile:      "",
		ConnectionString: "",
		KeySpace:         "flowg",
		GCInterval:       5 * time.Minute,
	}
}

// adapterOptions returns the options of the underlying FoundationDB adapter.
func (o Options) adapterOptions() foundation.AdapterOptions {
	return foundation.AdapterOptions{
		ClusterFile:      o.ClusterFile,
		ConnectionString: o.ConnectionString,
		KeySpace:         o.KeySpace,
		Namespace:        "log",
		GCInterval:       o.GCInterval,
	}
}

// Module returns an fx module providing a FoundationDB-backed
// [storage.LogStorage] configured by the [Options] provided to the graph. It
// also starts a background worker that periodically enforces each stream's
// retention budget.
func Module() fx.Option {
	return fx.Module(
		"storage.log",
		foundation.NewAdapter("storage.log", Options.adapterOptions),
		fx.Provide(func(lc fx.Lifecycle, d deps) storage.LogStorage {
			storage := log.NewStorage(d.Logger, d.Adapter, d.Options.BatchSize, d.Options.StreamConfigOverride)

			gc := actor.New(log.NewGarbageCollector(storage, d.Options.GCInterval))

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
