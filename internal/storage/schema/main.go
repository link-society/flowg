package schema

import (
	"context"
	"fmt"

	"bytes"
	"encoding/binary"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/pb"

	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/kvstore"
	"link-society.com/flowg/internal/utils/lww"
)

const CurrentVersion uint64 = 1

var versionKey = []byte("version")

func IsVersionKey(key []byte) bool {
	return bytes.Equal(key, versionKey)
}

const bitDelete = byte(1 << 0)

func IsTombstone(kv *pb.KV) bool {
	return len(kv.Meta) > 0 && kv.Meta[0]&bitDelete != 0
}

func ApplyEnvelope(txn *badger.Txn, key []byte, value []byte) (bool, error) {
	env, err := lww.Unmarshal(value)
	if err != nil {
		return false, err
	}

	return lww.Apply(txn, key, env)
}

type AppliedRecord struct {
	Key     []byte
	Deleted bool
	Origin  string
}

func ApplyRecord(txn *badger.Txn, key []byte, value []byte) (AppliedRecord, bool, error) {
	env, err := lww.Unmarshal(value)
	if err != nil {
		return AppliedRecord{}, false, err
	}

	applied, err := lww.Apply(txn, key, env)
	if err != nil {
		return AppliedRecord{}, false, err
	}
	if !applied {
		return AppliedRecord{}, false, nil
	}

	return AppliedRecord{
		Key:     append([]byte(nil), key...),
		Deleted: env.Deleted,
		Origin:  env.Timestamp.NodeID,
	}, true, nil
}

func MergeEnveloped(applied *[]AppliedRecord) func(txn *badger.Txn, kv *pb.KV) error {
	return func(txn *badger.Txn, kv *pb.KV) error {
		if IsVersionKey(kv.Key) {
			return nil
		}

		if IsTombstone(kv) {
			return nil
		}

		rec, ok, err := ApplyRecord(txn, kv.Key, kv.Value)
		if err != nil {
			return err
		}
		if ok {
			*applied = append(*applied, rec)
		}

		return nil
	}
}

func CollectGarbage(
	ctx context.Context,
	kvStore kvstore.Storage,
	grace time.Duration,
	prefixes [][]byte,
) (int, error) {
	before := hlc.Timestamp{WallTime: time.Now().Add(-grace).UnixNano()}

	purged := 0
	err := kvStore.Update(ctx, func(txn *badger.Txn) error {
		var err error
		purged, err = lww.CollectGarbage(txn, prefixes, before)
		return err
	})
	return purged, err
}

func Migrate(
	ctx context.Context,
	kvStore kvstore.Storage,
	clock *hlc.Clock,
	envelopePrefixes [][]byte,
) error {
	return kvStore.Update(ctx, func(txn *badger.Txn) error {
		version, err := readVersion(txn)
		if err != nil {
			return fmt.Errorf("could not read schema version: %w", err)
		}

		if version < 1 {
			if err := envelopeV0toV1(txn, clock.Now(), envelopePrefixes); err != nil {
				return fmt.Errorf("could not migrate schema from version 0 to version 1: %w", err)
			}

			version = 1
		}

		if err := writeVersion(txn, version); err != nil {
			return fmt.Errorf("could not write schema version: %w", err)
		}

		return nil
	})
}

func readVersion(txn *badger.Txn) (uint64, error) {
	item, err := txn.Get(versionKey)
	if err == badger.ErrKeyNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	var version uint64
	err = item.Value(func(val []byte) error {
		if len(val) != 8 {
			return fmt.Errorf("invalid schema version value of length %d", len(val))
		}
		version = binary.BigEndian.Uint64(val)
		return nil
	})
	if err != nil {
		return 0, err
	}

	return version, nil
}

func writeVersion(txn *badger.Txn, version uint64) error {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, version)
	return txn.Set(versionKey, buf)
}

func envelopeV0toV1(txn *badger.Txn, ts hlc.Timestamp, prefixes [][]byte) error {
	if len(prefixes) == 0 {
		prefixes = [][]byte{nil}
	}

	type entry struct {
		key     []byte
		payload []byte
	}

	pending := []entry{}

	for _, prefix := range prefixes {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		if prefix != nil {
			opts.Prefix = prefix
		}
		it := txn.NewIterator(opts)

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.KeyCopy(nil)
			if bytes.Equal(key, versionKey) {
				continue
			}

			payload, err := item.ValueCopy(nil)
			if err != nil {
				it.Close()
				return err
			}

			pending = append(pending, entry{key: key, payload: payload})
		}

		it.Close()
	}

	for _, e := range pending {
		env := lww.Envelope{Payload: e.payload}
		if err := txn.Set(e.key, env.Marshal()); err != nil {
			return fmt.Errorf("could not envelope key '%s': %w", e.key, err)
		}
	}

	return nil
}
