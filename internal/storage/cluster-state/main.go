package clusterstate

import (
	"context"

	"encoding/binary"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/utils/kvstore"
)

type Storage interface {
	FetchLocalState(ctx context.Context, nodeID string, endpoints []string) (*NodeState, error)
	UpdateLocalState(ctx context.Context, nodeID string, namespace string, since uint64) error
	GetLiveness(ctx context.Context, namespace string) (int64, error)
	SetLiveness(ctx context.Context, namespace string, unixNano int64) error
	ResetLocalState(ctx context.Context, namespace string) error
}

type Options struct {
	Directory string
	InMemory  bool
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
	kvOpts.InMemory = opts.InMemory

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
		prefix := []byte("lastsync:")
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := string(item.Key())

			rest := key[len(prefix):]
			sep := strings.LastIndex(rest, ":")
			if sep < 0 {
				continue
			}
			sourceNodeID := rest[:sep]
			namespace := rest[sep+1:]

			var since uint64
			err := item.Value(func(val []byte) error {
				since = binary.BigEndian.Uint64(val)
				return nil
			})
			if err != nil {
				return err
			}

			state.LastSync[sourceNodeID] = append(state.LastSync[sourceNodeID], NamespaceSyncState{
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

		item, err := txn.Get(key)
		switch err {
		case nil:
			var current uint64
			if err := item.Value(func(val []byte) error {
				current = binary.BigEndian.Uint64(val)
				return nil
			}); err != nil {
				return err
			}
			if current >= since {
				return nil
			}
		case badger.ErrKeyNotFound:
			// No watermark yet: fall through and record the first one.
		default:
			return err
		}

		return txn.Set(key, binary.BigEndian.AppendUint64(nil, since))
	})
}

func (s *storageImpl) GetLiveness(ctx context.Context, namespace string) (int64, error) {
	var ts int64
	err := s.kvStore.View(ctx, func(txn *badger.Txn) error {
		key := fmt.Appendf(nil, "liveness:%s", namespace)
		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			ts = int64(binary.BigEndian.Uint64(val))
			return nil
		})
	})
	if err != nil {
		return 0, err
	}
	return ts, nil
}

func (s *storageImpl) SetLiveness(ctx context.Context, namespace string, unixNano int64) error {
	return s.kvStore.Update(ctx, func(txn *badger.Txn) error {
		key := fmt.Appendf(nil, "liveness:%s", namespace)
		return txn.Set(key, binary.BigEndian.AppendUint64(nil, uint64(unixNano)))
	})
}

func (s *storageImpl) ResetLocalState(ctx context.Context, namespace string) error {
	return s.kvStore.Update(ctx, func(txn *badger.Txn) error {
		prefix := []byte("lastsync:")
		suffix := ":" + namespace

		it := txn.NewIterator(badger.DefaultIteratorOptions)
		var keys [][]byte
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			if strings.HasSuffix(string(key), suffix) {
				keys = append(keys, key)
			}
		}
		it.Close()

		for _, key := range keys {
			if err := txn.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}
