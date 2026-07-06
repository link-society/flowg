package server

import "go.uber.org/fx"

// FoundationDbStorageOptions implements the StorageOptions interface for the
// FoundationDB storage backend. It provides the fx modules that wire up the
// FoundationDB storage backend with its respective cluster file and key space.
type FoundationDbStorageOptions struct {
	ClusterFile string
	KeySpace    string
}

func (o FoundationDbStorageOptions) AuthModule() fx.Option {
	panic("not implemented") // TODO: Implement FoundationDB auth module
}

func (o FoundationDbStorageOptions) ConfigModule() fx.Option {
	panic("not implemented") // TODO: Implement FoundationDB config module
}

func (o FoundationDbStorageOptions) LogModule() fx.Option {
	panic("not implemented") // TODO: Implement FoundationDB log module
}
