package storage

import (
	"github.com/dgraph-io/badger/v3"
)

type fieldsIndex struct {
	txn    *badger.Txn
	stream string
}

func newFieldsIndex(txn *badger.Txn, stream string) *fieldsIndex {
	return &fieldsIndex{
		txn:    txn,
		stream: stream,
	}
}

func (index *fieldsIndex) Filter(filter Filter, allKeys []string) []string {
	allKeysMap := sliceToMap(allKeys)

	var evaluate func(f Filter) map[string]struct{}
	evaluate = func(f Filter) map[string]struct{} {
		switch f := f.(type) {
		case *AndFilter:
			keys := evaluate(f.Filters[0])

			for _, subFilter := range f.Filters[1:] {
				subKeys := evaluate(subFilter)
				keys = intersectKeysMap(keys, subKeys)
			}

			return keys

		case *OrFilter:
			keys := map[string]struct{}{}

			for _, subFilter := range f.Filters {
				subKeys := evaluate(subFilter)
				keys = unionKeysMap(keys, subKeys)
			}

			return keys

		case *NotFilter:
			keys := evaluate(f.Filter)
			return differenceKeysMap(allKeysMap, keys)

		case *FieldExact:
			fieldIndex := newFieldIndex(index.txn, index.stream, f.Field, f.Value)

			keys := map[string]struct{}{}

			fieldIndex.IterKeys(func(key string) {
				if _, found := allKeysMap[key]; found {
					keys[key] = struct{}{}
				}
			})

			return keys

		case *FieldIn:
			keys := map[string]struct{}{}

			for _, value := range f.Values {
				fieldIndex := newFieldIndex(index.txn, index.stream, f.Field, value)

				fieldIndex.IterKeys(func(key string) {
					if _, found := allKeysMap[key]; found {
						keys[key] = struct{}{}
					}
				})
			}

			return keys

		default:
			return map[string]struct{}{}
		}
	}

	keys := evaluate(filter)
	return mapToSlice(keys)
}
