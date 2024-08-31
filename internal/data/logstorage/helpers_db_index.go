package logstorage

import (
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v3"
)

func purgeEntryFromFieldIndex(txn *badger.Txn, stream string, key []byte) error {
	suffix := fmt.Sprintf(":%s", string(key))

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte(fmt.Sprintf("index:%s:field:", stream))
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		indexKey := string(item.Key())

		if strings.HasSuffix(indexKey, suffix) {
			if err := txn.Delete(item.KeyCopy(nil)); err != nil {
				return fmt.Errorf(
					"could not delete key '%s' from field index: %w",
					string(key), err,
				)
			}
		}
	}

	return nil
}
