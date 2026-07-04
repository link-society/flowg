package transactions

import (
	"fmt"

	"encoding/json"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"

	"link-society.com/flowg/internal/models"
)

// Ingest stores one log record and, in the same transaction, maintains the
// secondary data that makes it queryable. It writes the record under its
// entry key (packed via entrySS) with the stream's configured TTL, registers
// each of its fields in the stream's field set, and, for every field the stream
// is configured to index, adds an inverted-index key pointing back at the
// record.
//
// Unlike BadgerDB, FoundationDB has no per-key TTL, so time-based expiry is
// handled exclusively by the GC worker.
func Ingest(txn fdb.Transaction, stream string, logRecord *models.LogRecord, key []byte) error {
	val, err := json.Marshal(logRecord)
	if err != nil {
		return fmt.Errorf("could not marshal log entry: %w", err)
	}

	streamConfig, err := GetOrCreateStreamConfig(txn, stream)
	if err != nil {
		return err
	}

	txn.Set(fdb.Key(key), val)

	streamFieldSS := fieldSS.Sub(subspace.FromBytes([]byte(stream)))

	for field := range logRecord.Fields {
		fieldKey := streamFieldSS.Pack(tuple.Tuple{field})
		txn.Set(fieldKey, nil)

		if streamConfig.IsFieldIndexed(field) {
			value := logRecord.Fields[field]
			fieldIndex := newFieldIndex(txn, stream, field, value)
			if err := fieldIndex.AddKey(key); err != nil {
				return fmt.Errorf(
					"could not add field index '%s' of log entry to stream '%s': %w",
					field, stream, err,
				)
			}
		}
	}

	return nil
}
