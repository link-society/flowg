package cluster

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"link-society.com/flowg/internal/utils/kvstore"
)

type nodeState struct {
	NodeID   string                   `json:"node_id"`
	LastSync map[string]nodeSyncState `json:"last_sync"`
}

type nodeSyncState struct {
	Auth   uint64 `json:"auth"`
	Config uint64 `json:"config"`
	Log    uint64 `json:"log"`
}

func fetchLocalState(
	ctx context.Context,
	storage *kvstore.Storage,
	nodeID string,
	endpoints []string,
) (*nodeState, error) {
	state := &nodeState{
		NodeID:   nodeID,
		LastSync: map[string]nodeSyncState{},
	}

	err := storage.View(ctx, func(txn *badger.Txn) error {
		fetchLastSync := func(dbType string, endpointName string) (uint64, error) {
			key := fmt.Appendf(nil, "lastsync:%s:%s", dbType, endpointName)
			item, err := txn.Get(key)
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return 0, nil
				}

				return 0, err
			}

			var lastSync uint64
			item.Value(func(val []byte) error {
				lastSync = binary.BigEndian.Uint64(val)
				return nil
			})

			return lastSync, nil
		}

		for _, endpointName := range endpoints {
			lastAuthSync, err := fetchLastSync("auth", endpointName)
			if err != nil {
				return err
			}

			lastConfigSync, err := fetchLastSync("config", endpointName)
			if err != nil {
				return err
			}

			lastLogSync, err := fetchLastSync("log", endpointName)
			if err != nil {
				return err
			}

			state.LastSync[endpointName] = nodeSyncState{
				Auth:   lastAuthSync,
				Config: lastConfigSync,
				Log:    lastLogSync,
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch node state: %w", err)
	}

	return state, nil
}

func updateLocalState(
	ctx context.Context,
	storage *kvstore.Storage,
	nodeID string,
	dbType string,
	since uint64,
) error {
	key := fmt.Appendf(nil, "lastsync:%s:%s", dbType, nodeID)
	return storage.Update(ctx, func(txn *badger.Txn) error {
		return txn.Set(key, binary.BigEndian.AppendUint64(nil, since))
	})
}
