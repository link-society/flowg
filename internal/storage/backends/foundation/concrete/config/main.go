package config

import (
	"go.uber.org/fx"

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
func NewStorage(opts Options) fx.Option {
	adapterOpts := foundation.AdapterOptions{
		LogChannel:       "storage.config",
		ClusterFile:      opts.ClusterFile,
		ConnectionString: opts.ConnectionString,
		KeySpace:         opts.KeySpace,
		Namespace:        "config",
	}

	return fx.Module(
		"storage.config",
		foundation.NewAdapter(adapterOpts),
		fx.Provide(func(d deps) storage.ConfigStorage {
			return config.NewStorage(d.Adapter)
		}),
	)
}
