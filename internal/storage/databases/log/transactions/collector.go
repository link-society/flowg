package transactions

import (
	"fmt"
	"log/slog"

	"encoding/json"
	"time"

	"go.link-society.com/flowg/pkg/models"

	"go.link-society.com/flowg/internal/storage/generic/kv"
)

// Ingest stores one log record and, in the same transaction, maintains the
// secondary data that makes it queryable. It writes the record under its
// "entry:<stream>:..." key with the stream's retention TTL, registers each of
// its fields in the stream's field set ("stream:field:<stream>:<field>"), and,
// for every field the stream is configured to index, adds an inverted-index key
// pointing back at the record. When override is not nil, it is used as the
// stream's configuration (see [GetOrCreateStreamConfig]). Field values too
// large to be indexed are skipped, with a warning logged through logger.
func Ingest(
	logger *slog.Logger,
	txn kv.MutationTx,
	stream string,
	logRecord *models.LogRecord,
	key kv.Key,
	override *models.StreamConfig,
) error {
	val, err := json.Marshal(logRecord)
	if err != nil {
		return fmt.Errorf("could not marshal log entry: %w", err)
	}

	streamConfig, err := GetOrCreateStreamConfig(txn, stream, override)
	if err != nil {
		return err
	}

	if streamConfig.RetentionTime > 0 {
		ttl := time.Duration(streamConfig.RetentionTime) * time.Second
		err = txn.SetWithTTL(key, val, ttl)
	} else {
		err = txn.Set(key, val)
	}
	if err != nil {
		return fmt.Errorf(
			"could not add log entry '%s' to stream '%s': %w",
			key, stream, err,
		)
	}

	for field, value := range logRecord.Fields {
		fieldKey := kv.Key{"stream", "field", stream, field}
		if err := txn.Set(fieldKey, []byte{}); err != nil {
			return fmt.Errorf(
				"could not save field '%s' of log entry '%s' to stream '%s': %w",
				field, key, stream, err,
			)
		}

		if streamConfig.IsFieldIndexed(field) {
			fieldIndex := newFieldIndex(logger, txn, stream, field, value)
			if err := fieldIndex.AddKey(key, streamConfig.RetentionTime); err != nil {
				return fmt.Errorf(
					"could not add field index '%s' of log entry '%s' to stream '%s': %w",
					field, key, stream, err,
				)
			}
		}
	}

	return nil
}
