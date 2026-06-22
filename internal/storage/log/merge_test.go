package log

import (
	"encoding/json"
	"sync/atomic"
	"testing"

	"github.com/dgraph-io/badger/v4"
	badgerOptions "github.com/dgraph-io/badger/v4/options"
	"github.com/dgraph-io/badger/v4/pb"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/log/transactions"
	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/lww"
)

func newMergeTestDB(t *testing.T) *badger.DB {
	t.Helper()

	db, err := badger.Open(
		badger.DefaultOptions("").
			WithInMemory(true).
			WithLogger(nil).
			WithMemTableSize(8 << 20).
			WithCompression(badgerOptions.None).
			WithBlockCacheSize(0).
			WithIndexCacheSize(0),
	)
	if err != nil {
		t.Fatalf("open badger: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func mustUpdate(t *testing.T, db *badger.DB, fn func(txn *badger.Txn) error) {
	t.Helper()
	if err := db.Update(fn); err != nil {
		t.Fatalf("db.Update: %v", err)
	}
}

func keyExists(t *testing.T, db *badger.DB, key string) bool {
	t.Helper()
	found := false
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		found = true
		return nil
	})
	if err != nil {
		t.Fatalf("db.View: %v", err)
	}
	return found
}

func prefixExists(t *testing.T, db *badger.DB, prefix string) bool {
	t.Helper()
	found := false
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{Prefix: []byte(prefix)})
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			found = true
			return nil
		}
		return nil
	})
	if err != nil {
		t.Fatalf("db.View: %v", err)
	}
	return found
}

// TestMergeEntryDerivesFieldKeys verifies that merging a replicated log entry
// rebuilds the locally-derived field-presence and field-index keys. These keys
// are not replicated, so a peer must materialize them from the entry itself.
func TestMergeEntryDerivesFieldKeys(t *testing.T) {
	db := newMergeTestDB(t)

	ts := hlc.Timestamp{WallTime: 1000, NodeID: "a"}
	mustUpdate(t, db, func(txn *badger.Txn) error {
		return transactions.ConfigureStream(
			txn, "s",
			models.StreamConfig{IndexedFields: []string{"appname"}},
			ts,
		)
	})

	rec := models.NewLogRecord(map[string]string{"appname": "robot", "message": "hello"})
	val, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal record: %v", err)
	}
	key := []byte("entry:s:00000000000000000001:uuid-1")

	var changed atomic.Bool
	merge := mergeRecord(&changed)
	mustUpdate(t, db, func(txn *badger.Txn) error {
		return merge(txn, &pb.KV{Key: key, Value: val})
	})

	if !keyExists(t, db, "stream:field:s:appname") {
		t.Errorf("field key 'stream:field:s:appname' should be derived on entry merge")
	}
	if !keyExists(t, db, "stream:field:s:message") {
		t.Errorf("field key 'stream:field:s:message' should be derived on entry merge")
	}
	if !prefixExists(t, db, "index:s:field:appname:") {
		t.Errorf("index key for indexed field 'appname' should be derived on entry merge")
	}
	if prefixExists(t, db, "index:s:field:message:") {
		t.Errorf("field 'message' is not indexed; no index key should be derived")
	}
	if !changed.Load() {
		t.Errorf("merging a live entry should mark the store as changed")
	}
}

// TestMergeSkipsReplicatedDerivedTombstones is the core regression guard for the
// field-key resurrection bug: a peer that purged then re-ingested a stream would
// push stale field/index tombstones during anti-entropy. Applying those
// tombstones would delete the live derived keys on a healthy node, leaving
// listStreamFields permanently empty. Derived keys must never be applied from
// replication.
func TestMergeSkipsReplicatedDerivedTombstones(t *testing.T) {
	db := newMergeTestDB(t)

	mustUpdate(t, db, func(txn *badger.Txn) error {
		if err := txn.Set([]byte("stream:field:s:appname"), []byte{}); err != nil {
			return err
		}
		return txn.Set([]byte("index:s:field:appname:dmFs:entry:s:1:u"), []byte{})
	})

	var changed atomic.Bool
	merge := mergeRecord(&changed)

	// A peer pushes tombstones for the derived keys (meta has the badger delete bit).
	mustUpdate(t, db, func(txn *badger.Txn) error {
		if err := merge(txn, &pb.KV{Key: []byte("stream:field:s:appname"), Meta: []byte{1}}); err != nil {
			return err
		}
		return merge(txn, &pb.KV{Key: []byte("index:s:field:appname:dmFs:entry:s:1:u"), Meta: []byte{1}})
	})

	if !keyExists(t, db, "stream:field:s:appname") {
		t.Errorf("live field key must survive a replicated tombstone (resurrection bug)")
	}
	if !keyExists(t, db, "index:s:field:appname:dmFs:entry:s:1:u") {
		t.Errorf("live index key must survive a replicated tombstone (resurrection bug)")
	}
	if changed.Load() {
		t.Errorf("merging a replicated derived-key tombstone must be a no-op")
	}
}

// TestMergeConfigReindexesEntries verifies that when a stream config change is
// replicated, the peer rebuilds its local field indices for existing entries.
// Index keys are derived state and never replicated, so toggling a field's
// indexed status must trigger a local re-index on merge.
func TestMergeConfigReindexesEntries(t *testing.T) {
	db := newMergeTestDB(t)

	ts1 := hlc.Timestamp{WallTime: 1000, NodeID: "a"}
	mustUpdate(t, db, func(txn *badger.Txn) error {
		return transactions.ConfigureStream(
			txn, "s",
			models.StreamConfig{IndexedFields: []string{}},
			ts1,
		)
	})

	rec := models.NewLogRecord(map[string]string{"appname": "robot"})
	key := []byte("entry:s:00000000000000000001:uuid-1")
	mustUpdate(t, db, func(txn *badger.Txn) error {
		_, err := transactions.Ingest(txn, "s", rec, key, ts1)
		return err
	})

	if prefixExists(t, db, "index:s:field:appname:") {
		t.Fatalf("precondition failed: appname should not be indexed yet")
	}

	ts2 := hlc.Timestamp{WallTime: 2000, NodeID: "a"}
	cfgPayload, err := json.Marshal(models.StreamConfig{IndexedFields: []string{"appname"}})
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	cfgVal := lww.Envelope{Timestamp: ts2, Payload: cfgPayload}.Marshal()

	var changed atomic.Bool
	merge := mergeRecord(&changed)
	mustUpdate(t, db, func(txn *badger.Txn) error {
		return merge(txn, &pb.KV{Key: []byte("stream:config:s"), Value: cfgVal})
	})

	if !prefixExists(t, db, "index:s:field:appname:") {
		t.Errorf("merging a config that indexes 'appname' should backfill the index for existing entries")
	}
}
