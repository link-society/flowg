package transactions

import (
	"fmt"

	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/hlc"
)

func Ingest(txn *badger.Txn, stream string, logRecord *models.LogRecord, key []byte, ts hlc.Timestamp) (bool, error) {
	val, err := json.Marshal(logRecord)
	if err != nil {
		return false, fmt.Errorf("could not marshal log entry: %w", err)
	}

	streamConfig, created, err := GetOrCreateStreamConfig(txn, stream, ts)
	if err != nil {
		return false, err
	}

	entry := badger.NewEntry(key, val)
	if streamConfig.RetentionTime > 0 {
		entry = entry.WithTTL(time.Duration(streamConfig.RetentionTime) * time.Second)
	}

	if err := txn.SetEntry(entry); err != nil {
		return false, fmt.Errorf(
			"could not add log entry '%s' to stream '%s': %w",
			key, stream, err,
		)
	}

	for field, value := range logRecord.Fields {
		fieldKey := []byte(fmt.Sprintf("stream:field:%s:%s", stream, field))
		if err := txn.Set(fieldKey, []byte{}); err != nil {
			return false, fmt.Errorf(
				"could not save field '%s' of log entry '%s' to stream '%s': %w",
				field, key, stream, err,
			)
		}

		if streamConfig.IsFieldIndexed(field) {
			fieldIndex := newFieldIndex(txn, stream, field, value)
			if err := fieldIndex.AddKey(key, streamConfig.RetentionTime); err != nil {
				return false, fmt.Errorf(
					"could not add field index '%s' of log entry '%s' to stream '%s': %w",
					field, key, stream, err,
				)
			}
		}
	}

	return created, nil
}
