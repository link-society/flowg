package log

import (
	"context"
	"time"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/storage/backends/foundation"
	"link-society.com/flowg/internal/storage/databases/log"
	storage "link-society.com/flowg/internal/storage/interfaces"
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
}

type deps struct {
	fx.In

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

// NewStorage returns an fx module providing a FoundationDB-backed
// [storage.LogStorage] configured with the given options. It also starts a
// background worker that periodically enforces each stream's retention budget.
func NewStorage(opts Options) fx.Option {
	adapterOpts := foundation.AdapterOptions{
		LogChannel:       "storage.log",
		ClusterFile:      opts.ClusterFile,
		ConnectionString: opts.ConnectionString,
		KeySpace:         opts.KeySpace,
		Namespace:        "log",
		GCInterval:       opts.GCInterval,
	}

	return fx.Module(
		"storage.log",
		foundation.NewAdapter(adapterOpts),
		fx.Provide(func(lc fx.Lifecycle, d deps) storage.LogStorage {
			storage := log.NewStorage(d.Adapter, opts.BatchSize)

			gc := actor.New(log.NewGarbageCollector(storage, opts.GCInterval))

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
