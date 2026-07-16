package server

import (
	"go.uber.org/fx"

	"link-society.com/flowg/internal/storage/backends/foundation/concrete/auth"
	"link-society.com/flowg/internal/storage/backends/foundation/concrete/config"
	"link-society.com/flowg/internal/storage/backends/foundation/concrete/log"
)

// FoundationDbStorageOptions implements the StorageOptions interface for the
// FoundationDB storage backend. It provides the fx modules that wire up the
// FoundationDB storage backend with its respective cluster file and key space.
type FoundationDbStorageOptions struct {
	ClusterFile      string
	ConnectionString string
	KeySpace         string
}

func (o FoundationDbStorageOptions) AuthModule() fx.Option {
	storageOpts := auth.DefaultOptions()
	storageOpts.ClusterFile = o.ClusterFile
	storageOpts.ConnectionString = o.ConnectionString
	if o.KeySpace != "" {
		storageOpts.KeySpace = o.KeySpace
	}

	return auth.NewStorage(storageOpts)
}

func (o FoundationDbStorageOptions) ConfigModule() fx.Option {
	storageOpts := config.DefaultOptions()
	storageOpts.ClusterFile = o.ClusterFile
	storageOpts.ConnectionString = o.ConnectionString
	if o.KeySpace != "" {
		storageOpts.KeySpace = o.KeySpace
	}

	return config.NewStorage(storageOpts)
}

func (o FoundationDbStorageOptions) LogModule() fx.Option {
	storageOpts := log.DefaultOptions()
	storageOpts.ClusterFile = o.ClusterFile
	storageOpts.ConnectionString = o.ConnectionString
	if o.KeySpace != "" {
		storageOpts.KeySpace = o.KeySpace
	}

	return log.NewStorage(storageOpts)
}
