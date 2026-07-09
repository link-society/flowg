package interfaces

import (
	"context"
	"time"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/langs/filtering"
)

// LogStorage is the contract for persisting and querying FlowG's log records,
// organized into streams, along with each stream's configuration and field
// indices.
//
// It embeds [Streamable] so the whole store can be backed up and restored as a
// stream. Implementations live under internal/storage/backends.
type LogStorage interface {
	Streamable

	// ListStreamConfigs returns the configuration of every known stream, keyed
	// by stream name.
	ListStreamConfigs(ctx context.Context) (map[string]models.StreamConfig, error)
	// ListStreamFields returns the set of field names seen in the named stream.
	ListStreamFields(ctx context.Context, stream string) ([]string, error)
	// GetOrCreateStreamConfig returns the configuration of the named stream,
	// creating a default one if it does not exist yet.
	GetOrCreateStreamConfig(ctx context.Context, stream string) (models.StreamConfig, error)
	// ConfigureStream replaces the configuration of the named stream.
	ConfigureStream(ctx context.Context, stream string, config models.StreamConfig) error
	// DeleteStream removes the named stream and all of its records.
	DeleteStream(ctx context.Context, stream string) error
	// StreamUsage returns an estimate, in bytes, of the storage used by the
	// named stream.
	StreamUsage(ctx context.Context, stream string) (int64, error)

	// IndexField starts indexing the given field of the named stream.
	IndexField(ctx context.Context, stream, field string) error
	// UnindexField stops indexing the given field of the named stream.
	UnindexField(ctx context.Context, stream, field string) error
	// Distinct returns, for each indexed field of the named stream, the set of
	// distinct values observed.
	Distinct(ctx context.Context, stream string) (map[string][]string, error)

	// Ingest stores a log record in the named stream and returns the key it was
	// stored under.
	Ingest(ctx context.Context, stream string, logRecord *models.LogRecord) ([]byte, error)

	// FetchLogs returns the records of the named stream within the [from, to]
	// time range that satisfy the given filter and indexing constraints.
	FetchLogs(
		ctx context.Context,
		stream string,
		from, to time.Time,
		filter filtering.Filter,
		indexing map[string][]string,
	) ([]models.LogRecord, error)
}
