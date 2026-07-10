package transactions

import (
	"fmt"

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

// CollectGarbage evicts the oldest records of any stream that has grown past its
// configured size budget. For each over-budget stream it walks the entries from
// oldest to newest (the key order), deleting them and their inverted-index
// references until the stream fits within its retention size again.
func CollectGarbage(txn kv.MutationTx) error {
	streams, err := FetchStreamConfigs(txn)
	if err != nil {
		return err
	}

	for stream, config := range streams {
		retentionSize := config.RetentionSize * 1024 * 1024 // MB to bytes

		if retentionSize > 0 {
			streamSize := int64(0)

			for pair := range txn.IterPairs(kv.Key{"entry", stream}, kv.KeyRange{}) {
				streamSize += pair.EstimateSize()
			}

			if streamSize > retentionSize {
				for pair := range txn.IterPairs(kv.Key{"entry", stream}, kv.KeyRange{}) {
					streamSize -= pair.EstimateSize()

					key := pair.Key()

					if err := txn.Clear(key); err != nil {
						return fmt.Errorf(
							"could not delete key '%s' from stream '%s': %w",
							key, stream, err,
						)
					}

					purgeEntryFromFieldIndex(txn, stream, key)

					if streamSize <= retentionSize {
						break
					}
				}
			}
		}
	}

	return nil
}
