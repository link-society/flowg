package logstorage

import (
	"context"
	"fmt"
	"log/slog"

	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type CollectorSystem struct {
	storage *Storage
}

func NewCollectorSystem(storage *Storage) *CollectorSystem {
	return &CollectorSystem{storage: storage}
}

func (sys *CollectorSystem) Ingest(
	ctx context.Context,
	stream string,
	logEntry *LogEntry,
) ([]byte, error) {
	for {
		key := logEntry.NewDbKey(stream)
		val, err := json.Marshal(logEntry)
		if err != nil {
			slog.DebugContext(
				ctx,
				"Could not marshal log entry",
				"channel", "logstorage",
				"stream", stream,
				"error", err.Error(),
			)
			return nil, fmt.Errorf("could not marshal log entry: %w", err)
		}

		err = sys.storage.db.Update(func(txn *badger.Txn) error {
			slog.DebugContext(
				ctx,
				"Fetch stream configuration",
				"channel", "logstorage",
				"stream", stream,
			)

			streamConfig, err := getOrCreateStreamConfig(txn, stream)
			if err != nil {
				return err
			}

			slog.DebugContext(
				ctx,
				"Save log entry in BadgerDB",
				"channel", "logstorage",
				"stream", stream,
				"key", key,
			)

			entry := badger.NewEntry(key, val)
			if streamConfig.RetentionTime > 0 {
				entry = entry.WithTTL(time.Duration(streamConfig.RetentionTime) * time.Second)
			}

			if err := txn.SetEntry(entry); err != nil {
				return fmt.Errorf(
					"could not add log entry '%s' to stream '%s': %w",
					key, stream, err,
				)
			}

			for field, value := range logEntry.Fields {
				fieldKey := []byte(fmt.Sprintf("stream:field:%s:%s", stream, field))
				if err := txn.Set(fieldKey, []byte{}); err != nil {
					return fmt.Errorf(
						"could not save field '%s' of log entry '%s' to stream '%s': %w",
						field, key, stream, err,
					)
				}

				if streamConfig.IsFieldIndexed(field) {
					slog.DebugContext(
						ctx,
						"Save field index in BadgerDB",
						"channel", "logstorage",
						"stream", stream,
						"key", key,
						"field", field,
					)

					fieldIndex := newFieldIndex(txn, stream, field, value)
					if err := fieldIndex.AddKey(key, streamConfig.RetentionTime); err != nil {
						return fmt.Errorf(
							"could not add field index '%s' of log entry '%s' to stream '%s': %w",
							field, key, stream, err,
						)
					}
				}
			}

			return nil
		})

		switch err {
		case nil:
			return key, nil

		case badger.ErrConflict:
			slog.DebugContext(
				ctx,
				"Retry log entry ingestion",
				"channel", "logstorage",
				"stream", stream,
				"key", key,
			)
			continue

		default:
			return nil, err
		}
	}
}
