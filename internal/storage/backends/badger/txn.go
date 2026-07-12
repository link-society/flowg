package badger

import (
	"fmt"
	"iter"
	"time"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/storage/generic/kv"
)

// BadgerTx adapts a BadgerDB transaction to the [kv.QueryTx] and
// [kv.MutationTx] contracts.
type BadgerTx struct {
	concrete *badger.Txn
}

var _ kv.QueryTx = (*BadgerTx)(nil)
var _ kv.MutationTx = (*BadgerTx)(nil)

// Get implements [kv.QueryTx]. Only a missing key (badger.ErrKeyNotFound)
// yields a nil value; a key that exists always yields a non-nil value — an empty
// []byte for an empty-valued key — so callers can use a nil check to test for
// existence.
func (txn *BadgerTx) Get(key kv.Key) (kv.Value, error) {
	bkey := keyToBadger(key)
	if err := kv.CheckKeySize(len(bkey)); err != nil {
		return nil, err
	}

	item, err := txn.concrete.Get(bkey)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}

		return nil, err
	}

	content, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}
	if content == nil {
		content = []byte{}
	}
	val := kv.Value(content)

	return val, nil
}

// IterKeys implements [kv.QueryTx]. Values are not prefetched.
func (txn *BadgerTx) IterKeys(prefix kv.Key, keyRange kv.KeyRange) iter.Seq[kv.Key] {
	return func(yield func(kv.Key) bool) {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = fmt.Appendf(nil, "%s:", keyToBadger(prefix))
		it := txn.concrete.NewIterator(opts)
		defer it.Close()

		var (
			fromPrefix []byte
			toPrefix   string
		)

		if keyRange.From != nil {
			fromPrefix = fmt.Appendf(nil, "%s:", keyToBadger(keyRange.From))
		}
		if keyRange.To != nil {
			toPrefix = fmt.Sprintf("%s:", keyToBadger(keyRange.To))
		}

		if fromPrefix != nil {
			it.Seek(fromPrefix)
		} else {
			it.Rewind()
		}
		for it.Valid() {
			item := &badgerPair{concrete: it.Item()}

			if !yield(item.Key()) {
				return
			}

			if toPrefix != "" {
				key := string(item.concrete.Key())
				if key >= toPrefix {
					break
				}
			}

			it.Next()
		}
	}
}

// IterPairs implements [kv.QueryTx]. Values are prefetched.
func (txn *BadgerTx) IterPairs(prefix kv.Key, keyRange kv.KeyRange) iter.Seq[kv.Pair] {
	return func(yield func(kv.Pair) bool) {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.Prefix = fmt.Appendf(nil, "%s:", keyToBadger(prefix))
		it := txn.concrete.NewIterator(opts)
		defer it.Close()

		var (
			fromPrefix []byte
			toPrefix   string
		)

		if keyRange.From != nil {
			fromPrefix = fmt.Appendf(nil, "%s:", keyToBadger(keyRange.From))
		}
		if keyRange.To != nil {
			toPrefix = fmt.Sprintf("%s:", keyToBadger(keyRange.To))
		}

		if fromPrefix != nil {
			it.Seek(fromPrefix)
		} else {
			it.Rewind()
		}

		for it.Valid() {
			item := &badgerPair{concrete: it.Item()}

			if !yield(item) {
				return
			}

			if toPrefix != "" {
				key := string(item.concrete.Key())
				if key >= toPrefix {
					break
				}
			}

			it.Next()
		}
	}
}

// Set implements [kv.MutationTx].
func (txn *BadgerTx) Set(key kv.Key, value kv.Value) error {
	bkey := keyToBadger(key)
	if err := kv.CheckKeySize(len(bkey)); err != nil {
		return err
	}
	if err := kv.CheckValueSize(len(value)); err != nil {
		return err
	}
	return txn.concrete.Set(bkey, value)
}

// SetWithTTL implements [kv.MutationTx].
func (txn *BadgerTx) SetWithTTL(key kv.Key, value kv.Value, ttl time.Duration) error {
	bkey := keyToBadger(key)
	if err := kv.CheckKeySize(len(bkey)); err != nil {
		return err
	}
	if err := kv.CheckValueSize(len(value)); err != nil {
		return err
	}
	entry := badger.NewEntry(bkey, value).WithTTL(ttl)
	return txn.concrete.SetEntry(entry)
}

// Clear implements [kv.MutationTx].
func (txn *BadgerTx) Clear(key kv.Key) error {
	bkey := keyToBadger(key)
	if err := kv.CheckKeySize(len(bkey)); err != nil {
		return err
	}
	return txn.concrete.Delete(bkey)
}
