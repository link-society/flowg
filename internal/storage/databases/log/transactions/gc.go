package transactions

import (
	"encoding/json"
	"fmt"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"
)

// EstimateStorage sums estimated on-disk size of every entry in a stream
// ("entry:<stream>:*").
func EstimateStorage(txn kv.QueryTx, stream string) (int64, error) {
	storage := int64(0)

	for pair := range txn.IterPairs(kv.Key{"entry", stream}, kv.KeyRange{}) {
		storage += pair.EstimateSize()
	}

	return storage, nil
}

// EvictOldestBatch deletes up to limit of the oldest entries of a stream
// (walking entry keys in chronological order), stopping early once the freed
// bytes reach toFree. It also drops each deleted entry's inverted-index
// references. It returns the number of bytes freed and the number of entries
// deleted; a deleted count of zero means the stream has no more entries.
//
// Callers loop across successive transactions — subtracting the freed bytes
// from their running size estimate — so retention enforcement stays within the
// backend's transaction size and time limits.
func EvictOldestBatch(txn kv.MutationTx, stream string, toFree int64, limit int) (int64, int, error) {
	type victim struct {
		key    kv.Key
		record models.LogRecord
	}

	var (
		victims []victim
		freed   int64
	)

	for pair := range txn.IterPairs(kv.Key{"entry", stream}, kv.KeyRange{}) {
		if len(victims) >= limit || freed >= toFree {
			break
		}

		key := pair.Key()

		var record models.LogRecord
		if err := json.Unmarshal(pair.Value(), &record); err != nil {
			return freed, len(victims), fmt.Errorf(
				"could not unmarshal log entry '%s': %w", key, err,
			)
		}

		victims = append(victims, victim{key: key, record: record})
		freed += pair.EstimateSize()
	}

	for i := range victims {
		key := victims[i].key

		if err := txn.Clear(key); err != nil {
			return freed, len(victims), fmt.Errorf(
				"could not delete key '%s' from stream '%s': %w",
				key, stream, err,
			)
		}

		if err := purgeEntryFromFieldIndex(txn, stream, key, &victims[i].record); err != nil {
			return freed, len(victims), err
		}
	}

	return freed, len(victims), nil
}
