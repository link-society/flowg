package auth

import (
	"context"
	"fmt"

	"strings"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/utils/kvstore"
)

func migrateAlertScopes(ctx context.Context, kvStore kvstore.Storage) error {
	return kvStore.Update(ctx, func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("role:")
		it := txn.NewIterator(opts)

		roleNames := struct {
			readers []string
			writers []string
		}{}

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			parts := strings.Split(string(key[5:]), ":")
			roleName := parts[0]
			scopeName := parts[1]

			switch scopeName {
			case "read_alerts":
				roleNames.readers = append(roleNames.readers, roleName)
			case "write_alerts":
				roleNames.writers = append(roleNames.writers, roleName)
			}
		}

		it.Close()

		updates := []struct {
			scopeType string
			roles     []string
		}{
			{"read", roleNames.readers},
			{"write", roleNames.writers},
		}

		for _, update := range updates {
			for _, roleName := range update.roles {
				oldKey := fmt.Sprintf("role:%s:%s_alerts", roleName, update.scopeType)
				newKey := fmt.Sprintf("role:%s:%s_forwarders", roleName, update.scopeType)

				if err := txn.Set([]byte(newKey), []byte{}); err != nil {
					return fmt.Errorf("could not migrate old scope '%s' to new scope '%s': %w", oldKey, newKey, err)
				}

				if err := txn.Delete([]byte(oldKey)); err != nil {
					return fmt.Errorf("could not delete old scope '%s': %w", oldKey, err)
				}
			}
		}

		return nil
	})
}
