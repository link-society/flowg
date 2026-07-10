package badger

import (
	"strings"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/storage/generic/kv"
)

func keyToBadger(key kv.Key) []byte {
	return []byte(strings.Join(key, ":"))
}

type badgerPair struct {
	concrete *badger.Item
}

var _ kv.Pair = (*badgerPair)(nil)

func (p *badgerPair) Key() kv.Key {
	key := p.concrete.Key()
	parts := strings.Split(string(key), ":")
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
