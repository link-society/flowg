package auth

import (
	"context"
	"fmt"

	"strings"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/kvstore"
	"link-society.com/flowg/internal/utils/lww"
)

func migrateAlertScopes(ctx context.Context, kvStore kvstore.Storage, clock *hlc.Clock) error {
	return kvStore.Update(ctx, func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("role:")
		it := txn.NewIterator(opts)

		roleNames := struct {
			readers []string
			writers []string
		}{}

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			live := false
			err := item.Value(func(val []byte) error {
				env, err := lww.Unmarshal(val)
				if err != nil {
					return err
				}
				live = !env.Deleted
				return nil
			})
			if err != nil {
				it.Close()
				return err
			}
			if !live {
				continue
			}

			parts := strings.Split(string(item.Key()[5:]), ":")
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

				ts := clock.Now()

				if _, err := lww.Apply(txn, []byte(newKey), lww.Envelope{Timestamp: ts}); err != nil {
					return fmt.Errorf("could not migrate old scope '%s' to new scope '%s': %w", oldKey, newKey, err)
				}

				if _, err := lww.Apply(txn, []byte(oldKey), lww.Envelope{Timestamp: ts, Deleted: true}); err != nil {
					return fmt.Errorf("could not delete old scope '%s': %w", oldKey, err)
				}
			}
		}

		return nil
	})
}
