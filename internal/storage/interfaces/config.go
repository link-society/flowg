package interfaces

import (
	"context"

	"link-society.com/flowg/internal/models"
)

// ConfigStorage is the contract for persisting FlowG's configuration objects:
// transformers, pipelines, forwarders, and the global system configuration.
//
// It embeds [Streamable] so the whole store can be backed up and restored as a
// stream. Implementations live under internal/storage/backends.
type ConfigStorage interface {
	Streamable

	// ListTransformers returns the names of every stored transformer.
	ListTransformers(ctx context.Context) ([]string, error)
	// ReadTransformer returns the VRL source of the named transformer.
	ReadTransformer(ctx context.Context, name string) (string, error)
	// WriteTransformer creates or replaces the named transformer with the given
	// VRL source.
	WriteTransformer(ctx context.Context, name string, content string) error
	// DeleteTransformer removes the named transformer.
	DeleteTransformer(ctx context.Context, name string) error

	// ListPipelines returns the names of every stored pipeline.
	ListPipelines(ctx context.Context) ([]string, error)
	// ReadPipeline returns the named pipeline, migrating it to the latest flow
	// graph version if necessary.
	ReadPipeline(ctx context.Context, name string) (*models.FlowGraphV2, error)
	// WritePipeline creates or replaces the named pipeline from a flow graph.
	WritePipeline(ctx context.Context, name string, flow *models.FlowGraphV2) error
	// WriteRawPipeline creates or replaces the named pipeline from its raw
	// serialized form.
	WriteRawPipeline(ctx context.Context, name string, content string) error
	// DeletePipeline removes the named pipeline.
	DeletePipeline(ctx context.Context, name string) error

	// ListForwarders returns the names of every stored forwarder.
	ListForwarders(ctx context.Context) ([]string, error)
	// ReadForwarder returns the named forwarder, migrating it to the latest
	// version if necessary.
	ReadForwarder(ctx context.Context, name string) (*models.ForwarderV2, error)
	// WriteForwarder creates or replaces the named forwarder.
	WriteForwarder(ctx context.Context, name string, forwarder *models.ForwarderV2) error
	// DeleteForwarder removes the named forwarder.
	DeleteForwarder(ctx context.Context, name string) error

	// HasSystemConfig reports whether a system configuration has been stored.
	HasSystemConfig(ctx context.Context) (bool, error)
	// ReadSystemConfig returns the stored system configuration, or a zero value
	// if none has been written yet.
	ReadSystemConfig(ctx context.Context) (*models.SystemConfiguration, error)
	// WriteSystemConfig validates and stores the system configuration.
	WriteSystemConfig(ctx context.Context, config *models.SystemConfiguration) error
}
