package server

import (
	"go.uber.org/fx"

	foundationdb_auth "link-society.com/flowg/internal/storage/backends/foundationdb/concrete/auth"
	foundationdb_config "link-society.com/flowg/internal/storage/backends/foundationdb/concrete/config"
	foundationdb_log "link-society.com/flowg/internal/storage/backends/foundationdb/concrete/log"
)

// FoundationDbStorageOptions implements the StorageOptions interface for the
// FoundationDB storage backend. It provides the fx modules that wire up the
// three storage backends (auth, config, log) with their respective key prefixes.
type FoundationDbStorageOptions struct {
	ConnectionString string
	Prefix           []byte
}

func (o FoundationDbStorageOptions) AuthModule() fx.Option {
	storageOpts := foundationdb_auth.DefaultOptions()
	storageOpts.ConnectionString = o.ConnectionString
	storageOpts.Prefix = append(o.Prefix, []byte("auth/")...)

	return foundationdb_auth.NewStorage(storageOpts)
}

func (o FoundationDbStorageOptions) ConfigModule() fx.Option {
	storageOpts := foundationdb_config.DefaultOptions()
	storageOpts.ConnectionString = o.ConnectionString
	storageOpts.Prefix = append(o.Prefix, []byte("config/")...)

	return foundationdb_config.NewStorage(storageOpts)
}

func (o FoundationDbStorageOptions) LogModule() fx.Option {
	storageOpts := foundationdb_log.DefaultOptions()
	storageOpts.ConnectionString = o.ConnectionString
	storageOpts.Prefix = string(append(o.Prefix, []byte("log/")...))

	return foundationdb_log.NewStorage(storageOpts)
}
