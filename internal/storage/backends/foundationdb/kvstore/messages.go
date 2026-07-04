package kvstore

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
)

type message struct {
	replyTo chan<- error
	operation
}

type operation interface {
	Handle(db fdb.Database) error
}

type backupOperation struct {
	w      io.Writer
	since  uint64
	prefix []byte
}

type restoreOperation struct {
	r      io.Reader
	prefix []byte
}

type viewOperation struct {
	txnFn func(txn fdb.ReadTransaction) error
}

type updateOperation struct {
	txnFn func(txn fdb.Transaction) error
}

var _ operation = (*backupOperation)(nil)
var _ operation = (*restoreOperation)(nil)
var _ operation = (*viewOperation)(nil)
var _ operation = (*updateOperation)(nil)

type backupEntry struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Version uint64 `json:"version"`
}

func (op *backupOperation) Handle(db fdb.Database) error {
	_, err := db.ReadTransact(func(rtr fdb.ReadTransaction) (interface{}, error) {
		rv, err := rtr.GetReadVersion().Get()
		if err != nil {
			return nil, fmt.Errorf("failed to get read version: %w", err)
		}

		sub := subspace.FromBytes(op.prefix)
		begin, end := sub.FDBRangeKeys()
		kr := fdb.KeyRange{Begin: begin.FDBKey(), End: end.FDBKey()}

		ri := rtr.GetRange(kr, fdb.RangeOptions{})
		iter := ri.Iterator()
		enc := json.NewEncoder(op.w)

		for iter.Advance() {
			kv := iter.MustGet()

			entry := backupEntry{
				Key:     hex.EncodeToString(kv.Key),
				Value:   base64.StdEncoding.EncodeToString(kv.Value),
				Version: uint64(rv),
			}

			if err := enc.Encode(entry); err != nil {
				return nil, fmt.Errorf("failed to encode backup entry: %w", err)
			}
		}

		op.since = uint64(rv)

		return nil, nil
	})

	return err
}

func (op *restoreOperation) Handle(db fdb.Database) error {
	var entries []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	scanner := bufio.NewScanner(op.r)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024) // 1MB scan buffer

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := json.Unmarshal(line, &entry); err != nil {
			return fmt.Errorf("failed to decode backup entry: %w", err)
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading backup data: %w", err)
	}

	if len(entries) == 0 {
		return nil
	}

	const maxBatchSize = 100

	for i := 0; i < len(entries); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(entries) {
			end = len(entries)
		}
		batch := entries[i:end]

		_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			for _, entry := range batch {
				key, err := hex.DecodeString(entry.Key)
				if err != nil {
					return nil, fmt.Errorf("failed to decode key: %w", err)
				}

				value, err := base64.StdEncoding.DecodeString(entry.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to decode value: %w", err)
				}

				tr.Set(fdb.Key(key), value)
			}
			return nil, nil
		})
		if err != nil {
			return fmt.Errorf("failed to restore batch: %w", err)
		}
	}

	return nil
}

func (op *viewOperation) Handle(db fdb.Database) error {
	_, err := db.ReadTransact(func(rtr fdb.ReadTransaction) (interface{}, error) {
		return nil, op.txnFn(rtr)
	})
	return err
}

func (op *updateOperation) Handle(db fdb.Database) error {
	_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		return nil, op.txnFn(tr)
	})
	return err
}
