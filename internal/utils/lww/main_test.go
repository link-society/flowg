package lww_test

import (
	"bytes"
	"testing"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/lww"
)

func newDB(t *testing.T) *badger.DB {
	t.Helper()

	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true).WithLogger(nil),
	)
	if err != nil {
		t.Fatalf("open badger: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func ts(wall int64, logical uint32, node string) hlc.Timestamp {
	return hlc.Timestamp{WallTime: wall, Logical: logical, NodeID: node}
}

func TestMarshalRoundTrip(t *testing.T) {
	cases := []lww.Envelope{
		{Timestamp: ts(123, 4, "node-a"), Payload: []byte("hello")},
		{Timestamp: ts(0, 0, ""), Payload: nil},
		{Timestamp: ts(999, 1, "n"), Deleted: true, Payload: nil},
	}

	for _, want := range cases {
		got, err := lww.Unmarshal(want.Marshal())
		if err != nil {
			t.Fatalf("Unmarshal(%v) error = %v", want, err)
		}
		if !got.Timestamp.Equal(want.Timestamp) || got.Deleted != want.Deleted {
			t.Fatalf("got %+v; want %+v", got, want)
		}
		if !bytes.Equal(got.Payload, want.Payload) {
			t.Fatalf("payload got %q; want %q", got.Payload, want.Payload)
		}
	}
}

func TestUnmarshalRejectsGarbage(t *testing.T) {
	if _, err := lww.Unmarshal([]byte{0x01, 0x02}); err == nil {
		t.Fatal("expected error for short buffer")
	}
	if _, err := lww.Unmarshal(make([]byte, 16)); err == nil {
		t.Fatal("expected error for unsupported version 0")
	}
}

func apply(t *testing.T, db *badger.DB, key string, env lww.Envelope) bool {
	t.Helper()

	var applied bool
	err := db.Update(func(txn *badger.Txn) error {
		var aerr error
		applied, aerr = lww.Apply(txn, []byte(key), env)
		return aerr
	})
	if err != nil {
		t.Fatalf("Apply error = %v", err)
	}
	return applied
}

func read(t *testing.T, db *badger.DB, key string) (lww.Envelope, bool) {
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
		t.Fatalf("Read error = %v", err)
	}
	return env, found
}

func TestApplyWritesWhenAbsent(t *testing.T) {
	db := newDB(t)

	if !apply(t, db, "k", lww.Envelope{Timestamp: ts(100, 0, "a"), Payload: []byte("v1")}) {
		t.Fatal("expected first write to be applied")
	}

	env, found := read(t, db, "k")
	if !found || string(env.Payload) != "v1" {
		t.Fatalf("got (%q, %v); want (v1, true)", env.Payload, found)
	}
}

func TestApplyNewerWins(t *testing.T) {
	db := newDB(t)

	apply(t, db, "k", lww.Envelope{Timestamp: ts(100, 0, "a"), Payload: []byte("old")})

	if !apply(t, db, "k", lww.Envelope{Timestamp: ts(200, 0, "a"), Payload: []byte("new")}) {
		t.Fatal("expected newer write to be applied")
	}

	env, _ := read(t, db, "k")
	if string(env.Payload) != "new" {
		t.Fatalf("got %q; want new", env.Payload)
	}
}

func TestApplyOlderLoses(t *testing.T) {
	db := newDB(t)

	apply(t, db, "k", lww.Envelope{Timestamp: ts(200, 0, "a"), Payload: []byte("new")})

	if apply(t, db, "k", lww.Envelope{Timestamp: ts(100, 0, "a"), Payload: []byte("old")}) {
		t.Fatal("expected older write to lose")
	}

	env, _ := read(t, db, "k")
	if string(env.Payload) != "new" {
		t.Fatalf("got %q; want new", env.Payload)
	}
}

func TestApplyEqualIsIdempotent(t *testing.T) {
	db := newDB(t)

	apply(t, db, "k", lww.Envelope{Timestamp: ts(100, 1, "a"), Payload: []byte("v")})

	if apply(t, db, "k", lww.Envelope{Timestamp: ts(100, 1, "a"), Payload: []byte("v")}) {
		t.Fatal("expected equal timestamp to be a no-op")
	}
}

func TestApplyNodeIDTiebreak(t *testing.T) {
	db := newDB(t)

	apply(t, db, "k", lww.Envelope{Timestamp: ts(100, 0, "node-a"), Payload: []byte("a")})

	if !apply(t, db, "k", lww.Envelope{Timestamp: ts(100, 0, "node-b"), Payload: []byte("b")}) {
		t.Fatal("expected higher nodeID to win the tie")
	}

	env, _ := read(t, db, "k")
	if string(env.Payload) != "b" {
		t.Fatalf("got %q; want b", env.Payload)
	}
}

func TestTombstoneHidesValueAndBlocksOlderWrite(t *testing.T) {
	db := newDB(t)

	apply(t, db, "k", lww.Envelope{Timestamp: ts(100, 0, "a"), Payload: []byte("v")})

	if !apply(t, db, "k", lww.Envelope{Timestamp: ts(200, 0, "a"), Deleted: true}) {
		t.Fatal("expected delete to be applied")
	}

	if _, found := read(t, db, "k"); found {
		t.Fatal("expected tombstoned key to read as absent")
	}

	if apply(t, db, "k", lww.Envelope{Timestamp: ts(150, 0, "a"), Payload: []byte("zombie")}) {
		t.Fatal("expected write older than tombstone to lose")
	}

	if _, found := read(t, db, "k"); found {
		t.Fatal("tombstone should still hide the key after losing write")
	}
}

func TestWriteNewerThanTombstoneResurrects(t *testing.T) {
	db := newDB(t)

	apply(t, db, "k", lww.Envelope{Timestamp: ts(200, 0, "a"), Deleted: true})

	if !apply(t, db, "k", lww.Envelope{Timestamp: ts(300, 0, "a"), Payload: []byte("back")}) {
		t.Fatal("expected newer write to override tombstone")
	}

	env, found := read(t, db, "k")
	if !found || string(env.Payload) != "back" {
		t.Fatalf("got (%q, %v); want (back, true)", env.Payload, found)
	}
}

func TestCollectGarbage(t *testing.T) {
	db := newDB(t)

	apply(t, db, "old-tomb", lww.Envelope{Timestamp: ts(100, 0, "a"), Deleted: true})
	apply(t, db, "recent-tomb", lww.Envelope{Timestamp: ts(900, 0, "a"), Deleted: true})
	apply(t, db, "live", lww.Envelope{Timestamp: ts(100, 0, "a"), Payload: []byte("v")})

	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("version"), []byte("plain"))
	}); err != nil {
		t.Fatalf("seed plain key: %v", err)
	}

	var purged int
	if err := db.Update(func(txn *badger.Txn) error {
		var cerr error
		purged, cerr = lww.CollectGarbage(txn, nil, ts(500, 0, ""))
		return cerr
	}); err != nil {
		t.Fatalf("CollectGarbage error = %v", err)
	}

	if purged != 1 {
		t.Fatalf("purged = %d; want 1", purged)
	}

	assertAbsent := func(key string) {
		t.Helper()
		err := db.View(func(txn *badger.Txn) error {
			_, gerr := txn.Get([]byte(key))
			return gerr
		})
		if err != badger.ErrKeyNotFound {
			t.Fatalf("key %q: got err %v; want ErrKeyNotFound", key, err)
		}
	}

	assertPresent := func(key string) {
		t.Helper()
		err := db.View(func(txn *badger.Txn) error {
			_, gerr := txn.Get([]byte(key))
			return gerr
		})
		if err != nil {
			t.Fatalf("key %q: got err %v; want present", key, err)
		}
	}

	assertAbsent("old-tomb")
	assertPresent("recent-tomb")
	assertPresent("live")
	assertPresent("version")
}
