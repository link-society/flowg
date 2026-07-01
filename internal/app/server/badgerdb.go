package server

import (
	"go.uber.org/fx"

	"link-society.com/flowg/internal/storage/backends/badger/concrete/auth"
	"link-society.com/flowg/internal/storage/backends/badger/concrete/config"
	"link-society.com/flowg/internal/storage/backends/badger/concrete/log"
)

// BadgerDbStorageOptions implements the StorageOptions interface for the
// BadgerDB storage backend. It provides the fx modules that wire up the three
// storage backends (auth, config, log) with their respective on-disk directories.
type BadgerDbStorageOptions struct {
	AuthDir   string
	ConfigDir string
	LogDir    string
}

func (o BadgerDbStorageOptions) AuthModule() fx.Option {
	storageOpts := auth.DefaultOptions()
	storageOpts.Directory = o.AuthDir

	return auth.NewStorage(storageOpts)
}

func (o BadgerDbStorageOptions) ConfigModule() fx.Option {
	storageOpts := config.DefaultOptions()
	storageOpts.Directory = o.ConfigDir

	return config.NewStorage(storageOpts)
}

func (o BadgerDbStorageOptions) LogModule() fx.Option {
	storageOpts := log.DefaultOptions()
	storageOpts.Directory = o.LogDir

	return log.NewStorage(storageOpts)
}
