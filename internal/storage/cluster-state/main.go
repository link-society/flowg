package clusterstate

import (
	"context"

	"encoding/binary"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/utils/kvstore"
)

type Storage interface {
	FetchLocalState(ctx context.Context, nodeID string, endpoints []string) (*NodeState, error)
	UpdateLocalState(ctx context.Context, nodeID string, namespace string, since uint64) error
}

type Options struct {
	Directory string
}

type storageImpl struct {
	kvStore kvstore.Storage
}

type deps struct {
	fx.In

	S kvstore.Storage `name:"cluster.state"`
}

var _ Storage = (*storageImpl)(nil)

func DefaultOptions() Options {
	return Options{
		Directory: "",
	}
}

func NewStorage(opts Options) fx.Option {
	kvOpts := kvstore.DefaultOptions()
	kvOpts.LogChannel = "cluster.state"
	kvOpts.Directory = opts.Directory

	return fx.Module(
		"cluster.state",
		kvstore.NewStorage(kvOpts),
		fx.Provide(func(d deps) Storage {
			return &storageImpl{kvStore: d.S}
		}),
	)
}

func (s *storageImpl) FetchLocalState(ctx context.Context, nodeID string, endpoints []string) (*NodeState, error) {
	state := &NodeState{
		NodeID:   nodeID,
		LastSync: make(map[string][]NamespaceSyncState),
	}

	err := s.kvStore.View(ctx, func(txn *badger.Txn) error {
		prefix := fmt.Appendf(nil, "lastsync:%s:", nodeID)
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()
			namespace := string(key[len(prefix):])

			var since uint64
			err := item.Value(func(val []byte) error {
				since = binary.BigEndian.Uint64(val)
				return nil
			})
			if err != nil {
				return err
			}

			state.LastSync[namespace] = append(state.LastSync[namespace], NamespaceSyncState{
				Namespace: namespace,
				Since:     since,
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *storageImpl) UpdateLocalState(ctx context.Context, nodeID string, namespace string, since uint64) error {
	return s.kvStore.Update(ctx, func(txn *badger.Txn) error {
		key := fmt.Appendf(nil, "lastsync:%s:%s", nodeID, namespace)
		return txn.Set(key, binary.BigEndian.AppendUint64(nil, since))
	})
}
