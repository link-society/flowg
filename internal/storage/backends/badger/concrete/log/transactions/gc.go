package transactions

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

// EstimateStorage sums Badger's estimated on-disk size of every entry in a
// stream ("entry:<stream>:*").
func EstimateStorage(txn *badger.Txn, stream string) (int64, error) {
	storage := int64(0)

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = true
	opts.Prefix = []byte(fmt.Sprintf("entry:%s:", stream))
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		storage += item.EstimatedSize()
	}

	return storage, nil
}

// CollectGarbage evicts the oldest records of any stream that has grown past its
// configured size budget. For each over-budget stream it walks the entries from
// oldest to newest (the key order), deleting them and their inverted-index
// references until the stream fits within its retention size again.
func CollectGarbage(txn *badger.Txn) error {
	streams, err := FetchStreamConfigs(txn)
	if err != nil {
		return err
	}

	for stream, config := range streams {
		retentionSize := config.RetentionSize * 1024 * 1024 // MB to bytes

		if retentionSize > 0 {
			streamSize := int64(0)

			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = true
			opts.Prefix = []byte(fmt.Sprintf("entry:%s:", stream))
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				streamSize += item.EstimatedSize()
			}

			if streamSize > retentionSize {
				for it.Rewind(); it.Valid(); it.Next() {
					item := it.Item()
					streamSize -= item.EstimatedSize()

					key := item.KeyCopy(nil)

					if err := txn.Delete(key); err != nil {
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
