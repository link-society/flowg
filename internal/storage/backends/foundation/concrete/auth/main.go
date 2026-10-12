package auth

import (
	"go.uber.org/fx"

	"link-society.com/flowg/internal/storage/backends/foundation"
	"link-society.com/flowg/internal/storage/databases/auth"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// Options configures the FoundationDB-backed authentication storage.
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

	Adapter *foundation.FoundationAdapter `name:"storage.auth"`
}

// DefaultOptions returns the default [Options] for the authentication storage.
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
		Namespace:        "auth",
	}
}

// Module returns an fx module providing a FoundationDB-backed
// [storage.AuthStorage] configured by the [Options] provided to the graph.
func Module() fx.Option {
	return fx.Module(
		"storage.auth",
		foundation.NewAdapter("storage.auth", Options.adapterOptions),
		fx.Provide(func(d deps) storage.AuthStorage {
			return auth.NewStorage(d.Adapter)
		}),
	)
}
