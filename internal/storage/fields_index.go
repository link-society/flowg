package storage

import (
	"fmt"

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

func (index *fieldsIndex) Filter(filter Filter) ([]string, error) {
	var evaluate func(f Filter) (map[string]struct{}, error)
	evaluate = func(f Filter) (map[string]struct{}, error) {
		switch f := f.(type) {
		case *AndFilter:
			keys, err := evaluate(f.Filters[0])
			if err != nil {
				return nil, err
			}

			for _, subFilter := range f.Filters[1:] {
				subKeys, err := evaluate(subFilter)
				if err != nil {
					return nil, err
				}

				keys = intersectKeysMap(keys, subKeys)
			}

			return keys, nil

		case *OrFilter:
			keys := map[string]struct{}{}

			for _, subFilter := range f.Filters {
				subKeys, err := evaluate(subFilter)
				if err != nil {
					return nil, err
				}

				keys = unionKeysMap(keys, subKeys)
			}

			return keys, nil

		case *NotFilter:
			keys, err := evaluate(f.Filter)
			if err != nil {
				return nil, err
			}

			allKeys, err := index.getAllKeys()
			if err != nil {
				return nil, err
			}

			return differenceKeysMap(allKeys, keys), nil

		case *FieldExact:
			fieldIndex := newFieldIndex(index.txn, index.stream, f.Field, f.Value)
			keys, err := fieldIndex.GetKeys()
			if err != nil {
				return nil, err
			}

			return sliceToMap(keys), nil

		case *FieldIn:
			keys := map[string]struct{}{}

			for _, value := range f.Values {
				fieldIndex := newFieldIndex(index.txn, index.stream, f.Field, value)
				valueKeys, err := fieldIndex.GetKeys()
				if err != nil {
					return nil, err
				}

				for _, key := range valueKeys {
					keys[key] = struct{}{}
				}
			}

			return keys, nil

		default:
			return map[string]struct{}{}, nil
		}
	}

	keys, err := evaluate(filter)
	if err != nil {

		return nil, err
	}

	return mapToSlice(keys), err
}

func (index *fieldsIndex) GetKeys() ([]string, error) {
	keys, err := index.getAllKeys()
	if err != nil {
		return nil, err
	}

	return mapToSlice(keys), nil
}

func (index *fieldsIndex) getAllKeys() (map[string]struct{}, error) {
	keys := map[string]struct{}{}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte(fmt.Sprintf("entry:%s:", index.stream))
	it := index.txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := string(it.Item().KeyCopy(nil))
		keys[key] = struct{}{}
	}

	return keys, nil
}
