package auth

import (
	"fmt"

	"strings"

	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/utils/proctree"

	"github.com/dgraph-io/badger/v4"
	"link-society.com/flowg/internal/utils/kvstore"
)

type migratorProcH struct {
	kvStore kvstore.Storage
}

var _ proctree.ProcessHandler = (*migratorProcH)(nil)

func (p *migratorProcH) Init(ctx actor.Context) proctree.ProcessResult {
	if err := p.migrateAlertScopes(ctx); err != nil {
		return proctree.Terminate(err)
	}

	return proctree.Continue()
}

func (p *migratorProcH) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (p *migratorProcH) Terminate(ctx actor.Context, err error) error {
	return err
}

func (p *migratorProcH) migrateAlertScopes(ctx actor.Context) error {
	return p.kvStore.Update(ctx, func(txn *badger.Txn) error {
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
