package badger

import (
	"strings"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/storage/generic/kv"
)

// keySeparator delimits the segments of a composite [kv.Key] in BadgerDB's flat
// keyspace. It is the ASCII ESC control byte: a byte that FlowG disallows in the
// segment values it stores (stream names, field names, item names), so a segment
// can never contain it and the join/split round-trip is unambiguous. Change this
// single constant to re-encode every key.
const keySeparator = "\x1b"

func keyToBadger(key kv.Key) []byte {
	return []byte(strings.Join(key, keySeparator))
}

// keyToBadgerPrefix returns the encoded key followed by a trailing separator, so
// it matches a key's descendants at a segment boundary and never a sibling that
// merely shares a byte prefix (e.g. "role\x1bfoo" must not match "role\x1bfoobar").
func keyToBadgerPrefix(key kv.Key) []byte {
	return append(keyToBadger(key), keySeparator...)
}

type badgerPair struct {
	concrete *badger.Item
}

var _ kv.Pair = (*badgerPair)(nil)

func (p *badgerPair) Key() kv.Key {
	key := p.concrete.Key()
	parts := strings.Split(string(key), keySeparator)
	return kv.Key(parts)
}

func (p *badgerPair) Value() kv.Value {
	content, err := p.concrete.ValueCopy(nil)
	if err != nil {
		return nil
	}
	return kv.Value(content)
}

func (p *badgerPair) EstimateSize() int64 {
	return p.concrete.EstimatedSize()
}

func (p *badgerPair) ExpiresAt() uint64 {
	return p.concrete.ExpiresAt()
}
