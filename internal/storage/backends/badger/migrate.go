package badger

import (
	"bytes"
	"time"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/storage/generic/kv"
)

// legacyKeySeparator is the ':' byte that databases created before
// [keySeparator] used to join key segments.
const legacyKeySeparator = ":"

// keySeparatorMigrationMarker records that [migrateKeySeparator] has completed,
// so it never runs a second time. A second pass would be unsafe once a live key
// legitimately contains a ':' inside one of its segments (which the new,
// ESC-based encoding allows).
var keySeparatorMigrationMarker = keyToBadger(kv.Key{"__meta__", "migration", "key_separator"})

// migrateKeySeparator rewrites every key still joined with the legacy ':'
// separator so it uses [keySeparator] (ESC) instead. Replacing each ':' byte
// with the ESC byte is a pure re-encoding: it reproduces exactly the segments
// the old code parsed, so nothing but the on-disk delimiter changes.
//
// The rewrite streams through a [badger.WriteBatch], which flushes on its own,
// so it scales to arbitrarily large databases; every entry keeps its value and
// its remaining TTL. It is a no-op once the marker is present.
func migrateKeySeparator(db *badger.DB) error {
	migrated, err := keySeparatorMigrated(db)
	if err != nil {
		return err
	}
	if migrated {
		return nil
	}

	wb := db.NewWriteBatch()
	defer wb.Cancel()

	var (
		legacy = []byte(legacyKeySeparator)
		sep    = []byte(keySeparator)
		now    = time.Now()
	)

	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			key := item.KeyCopy(nil)
			if !bytes.Contains(key, legacy) {
				continue // already re-encoded
			}
			newKey := bytes.ReplaceAll(key, legacy, sep)

			// An entry already past its TTL is dropped rather than rewritten.
			var ttl time.Duration
			if expiresAt := item.ExpiresAt(); expiresAt != 0 {
				ttl = time.Unix(int64(expiresAt), 0).Sub(now)
				if ttl <= 0 {
					if err := wb.Delete(key); err != nil {
						return err
					}
					continue
				}
			}

			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			entry := badger.NewEntry(newKey, value)
			if ttl > 0 {
				entry = entry.WithTTL(ttl)
			}
			if err := wb.SetEntry(entry); err != nil {
				return err
			}
			if err := wb.Delete(key); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	if err := wb.Set(keySeparatorMigrationMarker, []byte{}); err != nil {
		return err
	}

	return wb.Flush()
}

// keySeparatorMigrated reports whether the key-separator migration marker is set.
func keySeparatorMigrated(db *badger.DB) (bool, error) {
	migrated := false
	err := db.View(func(txn *badger.Txn) error {
		switch _, err := txn.Get(keySeparatorMigrationMarker); err {
		case nil:
			migrated = true
			return nil
		case badger.ErrKeyNotFound:
			return nil
		default:
			return err
		}
	})
	return migrated, err
}
