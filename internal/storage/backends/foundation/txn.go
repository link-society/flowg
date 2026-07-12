package foundation

import (
	"iter"
	"time"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"

	"link-society.com/flowg/internal/storage/generic/kv"
)

// FoundationQueryTx adapts a read-only FoundationDB transaction to the
// [kv.QueryTx] contract.
//
// Every operation is scoped to sub, the subspace named after the storage
// (e.g. flowg/config). Keys handed to and returned from the transaction are
// always logical [kv.Key]s: the subspace prefix is applied on the way in and
// stripped on the way out, so consumers never observe it.
type FoundationQueryTx struct {
	concrete fdb.ReadTransaction
	sub      subspace.Subspace
}

// FoundationMutationTx adapts a read-write FoundationDB transaction to the
// [kv.MutationTx] contract.
//
// Like [FoundationQueryTx] every operation is scoped to sub, and reads observe
// the transaction's own pending writes.
type FoundationMutationTx struct {
	concrete fdb.Transaction
	sub      subspace.Subspace
}

var _ kv.QueryTx = (*FoundationQueryTx)(nil)
var _ kv.MutationTx = (*FoundationMutationTx)(nil)

// Get implements [kv.QueryTx]. A missing or expired key yields a nil value and
// no error.
func (txn *FoundationQueryTx) Get(key kv.Key) (kv.Value, error) {
	return txnGet(txn.concrete, txn.sub, key)
}

// IterKeys implements [kv.QueryTx].
func (txn *FoundationQueryTx) IterKeys(prefix kv.Key, keyRange kv.KeyRange) iter.Seq[kv.Key] {
	return txnIterKeys(txn.concrete, txn.sub, prefix, keyRange)
}

// IterPairs implements [kv.QueryTx].
func (txn *FoundationQueryTx) IterPairs(prefix kv.Key, keyRange kv.KeyRange) iter.Seq[kv.Pair] {
	return txnIterKvPairs(txn.concrete, txn.sub, prefix, keyRange)
}

// Get implements [kv.QueryTx]. A missing or expired key yields a nil value and
// no error.
func (txn *FoundationMutationTx) Get(key kv.Key) (kv.Value, error) {
	return txnGet(txn.concrete, txn.sub, key)
}

// IterKeys implements [kv.QueryTx].
func (txn *FoundationMutationTx) IterKeys(prefix kv.Key, keyRange kv.KeyRange) iter.Seq[kv.Key] {
	return txnIterKeys(txn.concrete, txn.sub, prefix, keyRange)
}

// IterPairs implements [kv.QueryTx].
func (txn *FoundationMutationTx) IterPairs(prefix kv.Key, keyRange kv.KeyRange) iter.Seq[kv.Pair] {
	return txnIterKvPairs(txn.concrete, txn.sub, prefix, keyRange)
}

// Set implements [kv.MutationTx].
func (txn *FoundationMutationTx) Set(key kv.Key, value kv.Value) error {
	fkey := keyToFdb(txn.sub, key)
	if err := kv.CheckKeySize(len(fkey)); err != nil {
		return err
	}
	if err := kv.CheckValueSize(expiryHeaderSize + len(value)); err != nil {
		return err
	}
	txn.concrete.Set(fkey, encodeValue(value, 0))
	return nil
}

// SetWithTTL implements [kv.MutationTx].
//
// FoundationDB has no native TTL, so the expiration is stored in the value
// envelope and enforced lazily when the pair is read.
func (txn *FoundationMutationTx) SetWithTTL(key kv.Key, value kv.Value, ttl time.Duration) error {
	fkey := keyToFdb(txn.sub, key)
	if err := kv.CheckKeySize(len(fkey)); err != nil {
		return err
	}
	if err := kv.CheckValueSize(expiryHeaderSize + len(value)); err != nil {
		return err
	}
	expiresAt := uint64(time.Now().Add(ttl).Unix())
	txn.concrete.Set(fkey, encodeValue(value, expiresAt))
	return nil
}

// Clear implements [kv.MutationTx].
func (txn *FoundationMutationTx) Clear(key kv.Key) error {
	fkey := keyToFdb(txn.sub, key)
	if err := kv.CheckKeySize(len(fkey)); err != nil {
		return err
	}
	txn.concrete.Clear(fkey)
	return nil
}

// txnGet reads a single key through read, honoring lazy TTL expiry.
func txnGet(read fdb.ReadTransaction, sub subspace.Subspace, key kv.Key) (kv.Value, error) {
	fkey := keyToFdb(sub, key)
	if err := kv.CheckKeySize(len(fkey)); err != nil {
		return nil, err
	}

	raw, err := read.Get(fkey).Get()
	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, nil
	}

	if expired(decodeExpiresAt(raw)) {
		return nil, nil
	}

	return decodeValue(raw), nil
}

// txnIterKeys iterates the keys contained in the prefix subspace.
func txnIterKeys(read fdb.ReadTransaction, sub subspace.Subspace, prefix kv.Key, keyRange kv.KeyRange) iter.Seq[kv.Key] {
	return func(yield func(kv.Key) bool) {
		for pair := range txnIterPairs(read, sub, prefix, keyRange) {
			if !yield(pair.Key()) {
				return
			}
		}
	}
}

// txnIterKvPairs iterates the pairs contained in the prefix subspace as
// [kv.Pair] values.
func txnIterKvPairs(read fdb.ReadTransaction, sub subspace.Subspace, prefix kv.Key, keyRange kv.KeyRange) iter.Seq[kv.Pair] {
	return func(yield func(kv.Pair) bool) {
		for pair := range txnIterPairs(read, sub, prefix, keyRange) {
			if !yield(pair) {
				return
			}
		}
	}
}

// txnIterPairs walks the key-value pairs contained in the prefix subspace,
// honoring the optional [kv.KeyRange] bounds and skipping expired entries.
//
// Iteration stops silently on error, mirroring the error-free iterator
// contract of [kv.QueryTx].
func txnIterPairs(read fdb.ReadTransaction, sub subspace.Subspace, prefix kv.Key, keyRange kv.KeyRange) iter.Seq[*foundationPair] {
	return func(yield func(*foundationPair) bool) {
		prefixSub := sub.Sub(keyToTuple(prefix)...)
		beginKey, endKey := prefixSub.FDBRangeKeys()

		var begin, end fdb.KeyConvertible = beginKey, endKey
		if keyRange.From != nil {
			begin = keyToFdb(sub, keyRange.From)
		}
		if keyRange.To != nil {
			end = keyToFdb(sub, keyRange.To)
		}

		r := fdb.KeyRange{Begin: begin, End: end}

		it := read.GetRange(r, fdb.RangeOptions{}).Iterator()
		for it.Advance() {
			concrete, err := it.Get()
			if err != nil {
				return
			}

			pair := &foundationPair{sub: sub, concrete: concrete}
			if expired(pair.ExpiresAt()) {
				continue
			}

			if !yield(pair) {
				return
			}
		}
	}
}
