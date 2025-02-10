package transactions

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

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
				streamSize += int64(item.EstimatedSize())
			}

			if streamSize > retentionSize {
				for it.Rewind(); it.Valid(); it.Next() {
					item := it.Item()
					streamSize -= int64(item.EstimatedSize())

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
