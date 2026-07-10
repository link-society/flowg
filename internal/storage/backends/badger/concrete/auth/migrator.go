package auth

import (
	"context"
	"fmt"

	"link-society.com/flowg/internal/storage/backends/badger"
	"link-society.com/flowg/internal/storage/generic/kv"
)

// migrateAlertScopes renames the legacy "alerts" permission to "forwarders".
// FlowG used to call forwarders "alerts", so older databases still hold
// "role:<name>:read_alerts" / "role:<name>:write_alerts" scope keys. This pass
// scans every role, and for each affected grant writes the equivalent
// "*_forwarders" key and deletes the obsolete "*_alerts" one.
func migrateAlertScopes(ctx context.Context, adapter *badger.BadgerAdapter) error {
	return adapter.Update(ctx, func(txn *badger.BadgerTx) error {
		roleNames := struct {
			readers []string
			writers []string
		}{}

		for key := range txn.IterKeys(kv.Key{"role"}, kv.KeyRange{}) {
			roleName := key[1]
			scopeName := key[2]

			switch scopeName {
			case "read_alerts":
				roleNames.readers = append(roleNames.readers, roleName)
			case "write_alerts":
				roleNames.writers = append(roleNames.writers, roleName)
			}
		}

		updates := []struct {
			scopeType string
			roles     []string
		}{
			{"read", roleNames.readers},
			{"write", roleNames.writers},
		}

		for _, update := range updates {
			for _, roleName := range update.roles {
				oldKey := kv.Key{"role", roleName, fmt.Sprintf("%s_alerts", update.scopeType)}
				newKey := kv.Key{"role", roleName, fmt.Sprintf("%s_forwarders", update.scopeType)}

				if err := txn.Set(newKey, []byte{}); err != nil {
					return fmt.Errorf("could not migrate old scope '%s' to new scope '%s': %w", oldKey, newKey, err)
				}

				if err := txn.Clear(oldKey); err != nil {
					return fmt.Errorf("could not delete old scope '%s': %w", oldKey, err)
				}
			}
		}

		return nil
	})
}
