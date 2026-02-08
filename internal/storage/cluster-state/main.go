package clusterstate

import (
	"go.uber.org/fx"

	"link-society.com/flowg/internal/utils/kvstore"
)

type Storage interface{}

type Options struct {
	Directory string
}

type storageImpl struct {
	kvStore kvstore.Storage
}

type deps struct {
	fx.In

	S kvstore.Storage `name:"cluster.state"`
}

var _ Storage = (*storageImpl)(nil)

func DefaultOptions() Options {
	return Options{
		Directory: "",
	}
}

func NewStorage(opts Options) fx.Option {
	kvOpts := kvstore.DefaultOptions()
	kvOpts.LogChannel = "cluster.state"
	kvOpts.Directory = opts.Directory

	return fx.Module(
		"storage.cluster.state",
		kvstore.NewStorage(kvOpts),
		fx.Provide(func(d deps) Storage {
			return &storageImpl{kvStore: d.S}
		}),
	)
}
