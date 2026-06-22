package kvstore

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/dgraph-io/badger/v4"
	badgerOptions "github.com/dgraph-io/badger/v4/options"
	"github.com/dgraph-io/badger/v4/pb"

	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/lww"
)

func newBadger(t *testing.T) *badger.DB {
	t.Helper()

	db, err := badger.Open(
		badger.DefaultOptions("").
			WithInMemory(true).
			WithLogger(nil).
			// An 8 MiB memtable with no compression or caches keeps
			// badger.Open cheap: the default 64 MiB memtable forces an
			// ~87 MiB arena allocation per open, which dominates test
			// runtime under CPU contention. 8 MiB is the floor because
			// in-memory mode pins the value threshold at 1 MiB.
			WithMemTableSize(8 << 20).
			WithCompression(badgerOptions.None).
			WithBlockCacheSize(0).
			WithIndexCacheSize(0),
	)
	if err != nil {
		t.Fatalf("failed to open badger: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return db
}

func putEnvelope(t *testing.T, db *badger.DB, key string, wall int64, payload string) {
	t.Helper()

	env := lww.Envelope{
		Timestamp: hlc.Timestamp{WallTime: wall, NodeID: "node"},
		Payload:   []byte(payload),
	}
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), env.Marshal())
	})
	if err != nil {
		t.Fatalf("failed to set %q: %v", key, err)
	}
}

func putRaw(t *testing.T, db *badger.DB, key string, value []byte) {
	t.Helper()

	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
	if err != nil {
		t.Fatalf("failed to set %q: %v", key, err)
	}
}

func readPayload(t *testing.T, db *badger.DB, key string) string {
	t.Helper()

	var payload string
	err := db.View(func(txn *badger.Txn) error {
		env, found, err := lww.Read(txn, []byte(key))
		if err != nil {
			return err
		}
		if !found {
			t.Fatalf("key %q not found", key)
		}
		payload = string(env.Payload)
		return nil
	})
	if err != nil {
		t.Fatalf("failed to read %q: %v", key, err)
	}

	return payload
}

func readRaw(t *testing.T, db *badger.DB, key string) []byte {
	t.Helper()

	var value []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		t.Fatalf("failed to read %q: %v", key, err)
	}

	return value
}

func mergeEnvelopeSkippingVersion(txn *badger.Txn, kv *pb.KV) error {
	if bytes.Equal(kv.Key, []byte("version")) {
		return nil
	}

	env, err := lww.Unmarshal(kv.Value)
	if err != nil {
		return err
	}

	_, err = lww.Apply(txn, kv.Key, env)
	return err
}

func TestMergeAppliesLWW(t *testing.T) {
	source := newBadger(t)
	putEnvelope(t, source, "k1", 10, "src-new")
	putEnvelope(t, source, "k2", 10, "src-old")

	versionBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(versionBuf, 99)
	putRaw(t, source, "version", versionBuf)

	var stream bytes.Buffer
	if _, err := source.Backup(&stream, 0); err != nil {
		t.Fatalf("failed to backup source: %v", err)
	}

	dest := newBadger(t)
	putEnvelope(t, dest, "k1", 5, "dst-stale")
	putEnvelope(t, dest, "k2", 20, "dst-fresh")

	localVersion := make([]byte, 8)
	binary.BigEndian.PutUint64(localVersion, 7)
	putRaw(t, dest, "version", localVersion)

	op := &mergeOperation{r: &stream, mergeFn: mergeEnvelopeSkippingVersion}
	if err := op.Handle(dest); err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	if got := readPayload(t, dest, "k1"); got != "src-new" {
		t.Errorf("k1: expected source to win (src-new), got %q", got)
	}

	if got := readPayload(t, dest, "k2"); got != "dst-fresh" {
		t.Errorf("k2: expected dest to win (dst-fresh), got %q", got)
	}

	if got := binary.BigEndian.Uint64(readRaw(t, dest, "version")); got != 7 {
		t.Errorf("version: expected local 7 to be preserved, got %d", got)
	}
}
