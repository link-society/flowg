package foundation

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"

	"link-society.com/flowg/internal/storage/generic/kv"
)

// expiryHeaderSize is the number of leading bytes reserved in every stored
// value to carry its expiration timestamp.
//
// FoundationDB has no native TTL support, so the expiration is embedded in the
// value envelope and enforced lazily on read.
const expiryHeaderSize = 8

// keyToTuple converts a composite [kv.Key] into a FoundationDB [tuple.Tuple].
//
// Each segment is a string, matching the [kv.Key] element type.
func keyToTuple(key kv.Key) tuple.Tuple {
	t := make(tuple.Tuple, len(key))
	for i, segment := range key {
		t[i] = segment
	}
	return t
}

// tupleToKey converts a FoundationDB [tuple.Tuple] back into a composite
// [kv.Key].
func tupleToKey(t tuple.Tuple) kv.Key {
	key := make(kv.Key, len(t))
	for i, element := range t {
		if s, ok := element.(string); ok {
			key[i] = s
		} else {
			key[i] = fmt.Sprintf("%v", element)
		}
	}
	return key
}

// keyToFdb packs a [kv.Key] into an [fdb.Key] scoped to the given subspace.
//
// The subspace prefix (e.g. flowg/config) is prepended by [subspace.Subspace.Pack]
// and must never leak back to consumers of the [kv.Adapter].
func keyToFdb(sub subspace.Subspace, key kv.Key) fdb.Key {
	return sub.Pack(keyToTuple(key))
}

// keyFromFdb unpacks an [fdb.Key] back into a [kv.Key], stripping the subspace
// prefix so consumers only ever see their own logical keys.
func keyFromFdb(sub subspace.Subspace, fdbKey fdb.Key) (kv.Key, error) {
	t, err := sub.Unpack(fdbKey)
	if err != nil {
		return nil, err
	}
	return tupleToKey(t), nil
}

// encodeValue wraps a payload with its expiration timestamp (unix seconds, 0
// meaning no expiration).
func encodeValue(value kv.Value, expiresAt uint64) []byte {
	buf := make([]byte, expiryHeaderSize+len(value))
	binary.BigEndian.PutUint64(buf[:expiryHeaderSize], expiresAt)
	copy(buf[expiryHeaderSize:], value)
	return buf
}

// decodeExpiresAt extracts the expiration timestamp from a stored value
// envelope. A malformed envelope is treated as non-expiring.
func decodeExpiresAt(raw []byte) uint64 {
	if len(raw) < expiryHeaderSize {
		return 0
	}
	return binary.BigEndian.Uint64(raw[:expiryHeaderSize])
}

// decodeValue extracts the payload from a stored value envelope.
func decodeValue(raw []byte) kv.Value {
	if len(raw) < expiryHeaderSize {
		return nil
	}
	return kv.Value(raw[expiryHeaderSize:])
}

// expired reports whether the given expiration timestamp has passed.
func expired(expiresAt uint64) bool {
	return expiresAt != 0 && uint64(time.Now().Unix()) >= expiresAt
}

// foundationPair adapts a FoundationDB [fdb.KeyValue] to the [kv.Pair] contract.
type foundationPair struct {
	sub      subspace.Subspace
	concrete fdb.KeyValue
}

var _ kv.Pair = (*foundationPair)(nil)

// Key implements [kv.Pair], returning the logical key with the subspace prefix
// removed.
func (p *foundationPair) Key() kv.Key {
	key, err := keyFromFdb(p.sub, p.concrete.Key)
	if err != nil {
		return nil
	}
	return key
}

// Value implements [kv.Pair], returning the payload without its expiry header.
func (p *foundationPair) Value() kv.Value {
	return decodeValue(p.concrete.Value)
}

// EstimateSize implements [kv.Pair].
func (p *foundationPair) EstimateSize() int64 {
	return int64(len(p.concrete.Key) + len(p.concrete.Value))
}

// ExpiresAt implements [kv.Pair], returning the embedded expiration timestamp
// (0 when the pair does not expire).
func (p *foundationPair) ExpiresAt() uint64 {
	return decodeExpiresAt(p.concrete.Value)
}
