package schema

import (
	"bytes"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/pb"

	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/lww"
)

func envTS(wall int64, node string) hlc.Timestamp {
	return hlc.Timestamp{WallTime: wall, NodeID: node}
}

func envelope(wall int64, node, payload string) []byte {
	return lww.Envelope{Timestamp: envTS(wall, node), Payload: []byte(payload)}.Marshal()
}

func tombstone(wall int64, node string) []byte {
	return lww.Envelope{Timestamp: envTS(wall, node), Deleted: true}.Marshal()
}

// putRaw writes a value directly, bypassing LWW, to seed a node's local state.
func putRaw(t *testing.T, db *badger.DB, key string, value []byte) {
	t.Helper()
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
	if err != nil {
		t.Fatalf("seed %q: %v", key, err)
	}
}

// dumpStream snapshots every stored key/value (including tombstone envelopes) as
// the anti-entropy backup stream would.
func dumpStream(t *testing.T, db *badger.DB) []*pb.KV {
	t.Helper()
	var out []*pb.KV
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			out = append(out, &pb.KV{Key: item.KeyCopy(nil), Value: value})
		}
		return nil
	})
	if err != nil {
		t.Fatalf("dump: %v", err)
	}
	return out
}

// mergeStream applies a backup stream into db through the production merge
// function and returns the records it actually applied.
func mergeStream(t *testing.T, db *badger.DB, stream []*pb.KV) []AppliedRecord {
	t.Helper()
	var applied []AppliedRecord
	fn := MergeEnveloped(&applied)
	err := db.Update(func(txn *badger.Txn) error {
		for _, kv := range stream {
			if err := fn(txn, kv); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	return applied
}

func readLive(t *testing.T, db *badger.DB, key string) (lww.Envelope, bool) {
	t.Helper()
	var (
		env   lww.Envelope
		found bool
	)
	err := db.View(func(txn *badger.Txn) error {
		var rerr error
		env, found, rerr = lww.Read(txn, []byte(key))
		return rerr
	})
	if err != nil {
		t.Fatalf("read %q: %v", key, err)
	}
	return env, found
}

func TestApplyRecordReportsOriginAndDeleted(t *testing.T) {
	db := newDB(t)

	var (
		rec AppliedRecord
		ok  bool
	)
	err := db.Update(func(txn *badger.Txn) error {
		var aerr error
		rec, ok, aerr = ApplyRecord(txn, []byte("k"), envelope(100, "node-a", "v"))
		return aerr
	})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	if !ok {
		t.Fatal("expected first write to be applied")
	}
	if rec.Origin != "node-a" {
		t.Errorf("origin: got %q want node-a", rec.Origin)
	}
	if rec.Deleted {
		t.Error("expected Deleted=false for a write")
	}
	if !bytes.Equal(rec.Key, []byte("k")) {
		t.Errorf("key: got %q want k", rec.Key)
	}
}

func TestApplyRecordTombstoneReportsDeleted(t *testing.T) {
	db := newDB(t)

	var rec AppliedRecord
	err := db.Update(func(txn *badger.Txn) error {
		var aerr error
		rec, _, aerr = ApplyRecord(txn, []byte("k"), tombstone(100, "node-b"))
		return aerr
	})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	if !rec.Deleted {
		t.Error("expected Deleted=true for a tombstone")
	}
	if rec.Origin != "node-b" {
		t.Errorf("origin: got %q want node-b", rec.Origin)
	}
}

func TestApplyRecordNotAppliedWhenOlder(t *testing.T) {
	db := newDB(t)

	putRaw(t, db, "k", envelope(200, "node-a", "winner"))

	var ok bool
	err := db.Update(func(txn *badger.Txn) error {
		var aerr error
		_, ok, aerr = ApplyRecord(txn, []byte("k"), envelope(100, "node-a", "loser"))
		return aerr
	})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	if ok {
		t.Fatal("expected an older record to be rejected")
	}
}

func TestMergeEnvelopedSkipsVersionKey(t *testing.T) {
	db := newDB(t)

	stream := []*pb.KV{
		{Key: versionKey, Value: []byte("anything")},
		{Key: []byte("k"), Value: envelope(100, "node-a", "v")},
	}
	applied := mergeStream(t, db, stream)

	if len(applied) != 1 {
		t.Fatalf("expected 1 applied record, got %d", len(applied))
	}
	if !bytes.Equal(applied[0].Key, []byte("k")) {
		t.Errorf("applied key: got %q want k", applied[0].Key)
	}

	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(versionKey)
		return err
	})
	if err != badger.ErrKeyNotFound {
		t.Errorf("expected version key to remain untouched, got err=%v", err)
	}
}

// TestBidirectionalConvergence proves that two nodes holding conflicting writes
// for the same keys converge to identical state after a round of anti-entropy in
// both directions, with the LWW winner chosen consistently on both sides.
func TestBidirectionalConvergence(t *testing.T) {
	nodeA := newDB(t)
	nodeB := newDB(t)

	// k1: A is newer; k2: B is newer; k3: only on A; k4: only on B.
	putRaw(t, nodeA, "k1", envelope(200, "a", "a-wins"))
	putRaw(t, nodeA, "k2", envelope(100, "a", "a-stale"))
	putRaw(t, nodeA, "k3", envelope(150, "a", "a-only"))

	putRaw(t, nodeB, "k1", envelope(100, "b", "b-stale"))
	putRaw(t, nodeB, "k2", envelope(200, "b", "b-wins"))
	putRaw(t, nodeB, "k4", envelope(150, "b", "b-only"))

	// Round 1: A pushes its full state to B.
	mergeStream(t, nodeB, dumpStream(t, nodeA))
	// Round 2: B pushes its (now-converged) full state back to A.
	mergeStream(t, nodeA, dumpStream(t, nodeB))

	want := map[string]string{
		"k1": "a-wins",
		"k2": "b-wins",
		"k3": "a-only",
		"k4": "b-only",
	}

	for _, node := range []struct {
		name string
		db   *badger.DB
	}{{"nodeA", nodeA}, {"nodeB", nodeB}} {
		for key, expected := range want {
			env, found := readLive(t, node.db, key)
			if !found {
				t.Errorf("%s: key %q missing after convergence", node.name, key)
				continue
			}
			if string(env.Payload) != expected {
				t.Errorf("%s: key %q = %q; want %q", node.name, key, env.Payload, expected)
			}
		}
	}
}

// TestTombstonePropagation proves a delete on one node propagates and wins on the
// other through the merge stream.
func TestTombstonePropagation(t *testing.T) {
	nodeA := newDB(t)
	nodeB := newDB(t)

	// Both nodes initially agree on the key.
	putRaw(t, nodeA, "k", envelope(100, "a", "v"))
	putRaw(t, nodeB, "k", envelope(100, "a", "v"))

	// A deletes it with a newer timestamp.
	putRaw(t, nodeA, "k", tombstone(200, "a"))

	applied := mergeStream(t, nodeB, dumpStream(t, nodeA))

	if _, found := readLive(t, nodeB, "k"); found {
		t.Fatal("expected key to be tombstoned on nodeB")
	}

	var sawTombstone bool
	for _, rec := range applied {
		if bytes.Equal(rec.Key, []byte("k")) && rec.Deleted {
			sawTombstone = true
		}
	}
	if !sawTombstone {
		t.Error("expected an applied tombstone record for key k")
	}
}

// TestMergeIsIdempotent proves replaying the same stream applies nothing the
// second time, so no spurious change events (and therefore no rebroadcasts) are
// produced by a redundant sync.
func TestMergeIsIdempotent(t *testing.T) {
	nodeA := newDB(t)
	nodeB := newDB(t)

	putRaw(t, nodeA, "k1", envelope(200, "a", "v1"))
	putRaw(t, nodeA, "k2", tombstone(200, "a"))

	stream := dumpStream(t, nodeA)

	first := mergeStream(t, nodeB, stream)
	if len(first) != 2 {
		t.Fatalf("first merge: expected 2 applied, got %d", len(first))
	}

	second := mergeStream(t, nodeB, stream)
	if len(second) != 0 {
		t.Fatalf("second merge: expected 0 applied (idempotent), got %d", len(second))
	}
}
