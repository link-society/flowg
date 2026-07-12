package badger

import (
	"testing"

	"errors"

	"slices"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/storage/generic/kv"
)

func newTestTx(t *testing.T) *BadgerTx {
	t.Helper()

	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true).WithLoggingLevel(badger.ERROR),
	)
	if err != nil {
		t.Fatalf("failed to open in-memory badger: %v", err)
	}

	txn := db.NewTransaction(true)
	t.Cleanup(func() {
		txn.Discard()
		db.Close()
	})

	return &BadgerTx{concrete: txn}
}

func TestSetRejectsOversizedValue(t *testing.T) {
	txn := newTestTx(t)

	key := kv.Key{"big"}
	value := kv.Value(strings.Repeat("x", kv.MaxValueSize+1))

	if err := txn.Set(key, value); !errors.Is(err, kv.ErrValueTooLarge) {
		t.Fatalf("expected ErrValueTooLarge, got %v", err)
	}

	ttlErr := txn.SetWithTTL(key, value, time.Minute)
	if !errors.Is(ttlErr, kv.ErrValueTooLarge) {
		t.Fatalf("expected ErrValueTooLarge from SetWithTTL, got %v", ttlErr)
	}
}

func TestSetRejectsOversizedKey(t *testing.T) {
	txn := newTestTx(t)

	key := kv.Key{strings.Repeat("k", kv.MaxKeySize+1)}
	value := kv.Value("ok")

	if err := txn.Set(key, value); !errors.Is(err, kv.ErrKeyTooLarge) {
		t.Fatalf("expected ErrKeyTooLarge, got %v", err)
	}

	ttlErr := txn.SetWithTTL(key, value, time.Minute)
	if !errors.Is(ttlErr, kv.ErrKeyTooLarge) {
		t.Fatalf("expected ErrKeyTooLarge from SetWithTTL, got %v", ttlErr)
	}
}

func TestSetAcceptsWithinLimits(t *testing.T) {
	txn := newTestTx(t)

	key := kv.Key{"ok"}
	value := kv.Value(strings.Repeat("x", kv.MaxValueSize))

	if err := txn.Set(key, value); err != nil {
		t.Fatalf("expected value at MaxValueSize to be accepted, got %v", err)
	}

	got, err := txn.Get(key)
	if err != nil {
		t.Fatalf("failed to read back value: %v", err)
	}
	if len(got) != kv.MaxValueSize {
		t.Fatalf("expected %d bytes, got %d", kv.MaxValueSize, len(got))
	}
}

// seedEntries writes entry-like keys ("entry:s:<ms>:u") for the given
// timestamp segments.
func seedEntries(t *testing.T, txn *BadgerTx, timestamps ...string) {
	t.Helper()
	for _, ms := range timestamps {
		if err := txn.Set(kv.Key{"entry", "s", ms, "u"}, kv.Value("x")); err != nil {
			t.Fatalf("failed to seed entry %q: %v", ms, err)
		}
	}
}

// iterKeyTimestamps collects the timestamp segment of every key IterKeys yields.
func iterKeyTimestamps(txn *BadgerTx, keyRange kv.KeyRange) []string {
	var got []string
	for key := range txn.IterKeys(kv.Key{"entry", "s"}, keyRange) {
		got = append(got, key[2])
	}
	return got
}

// TestIterKeysInclusiveTo checks that To is inclusive: the whole To subtree is
// returned (every record at the boundary timestamp) and nothing past it. This is
// the record that used to leak out inconsistently between backends.
func TestIterKeysInclusiveTo(t *testing.T) {
	txn := newTestTx(t)

	// Two records share the boundary timestamp 003; both must be returned.
	for _, key := range []kv.Key{
		{"entry", "s", "001", "a"},
		{"entry", "s", "002", "a"},
		{"entry", "s", "003", "a"},
		{"entry", "s", "003", "b"},
		{"entry", "s", "004", "a"},
	} {
		if err := txn.Set(key, kv.Value("x")); err != nil {
			t.Fatalf("failed to seed %v: %v", key, err)
		}
	}

	got := iterKeyTimestamps(txn, kv.KeyRange{
		From: kv.Key{"entry", "s", "002"},
		To:   kv.Key{"entry", "s", "003"},
	})

	want := []string{"002", "003", "003"} // both 003 records, never 004
	if !slices.Equal(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// TestIterKeysInclusiveFrom checks that the From bound is inclusive, matching the
// FoundationDB backend.
func TestIterKeysInclusiveFrom(t *testing.T) {
	txn := newTestTx(t)
	seedEntries(t, txn, "001", "002", "003")

	got := iterKeyTimestamps(txn, kv.KeyRange{From: kv.Key{"entry", "s", "002"}})

	want := []string{"002", "003"}
	if !slices.Equal(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// TestIterKeysToBoundSegmentBoundary guards the trailing-':' segment boundary:
// an inclusive To of {"role","foo"} returns the "role:foo" subtree, but a sibling
// like "role:foobar" (which shares the "foo" byte prefix) sorts past the boundary
// and never leaks in.
func TestIterKeysToBoundSegmentBoundary(t *testing.T) {
	txn := newTestTx(t)

	for _, key := range []kv.Key{
		{"role", "bar", "perm"},
		{"role", "foo", "perm"},
		{"role", "foobar", "perm"},
	} {
		if err := txn.Set(key, kv.Value("x")); err != nil {
			t.Fatalf("failed to seed %v: %v", key, err)
		}
	}

	// [role:bar, role:foo]: "bar" and the "foo" subtree; "foobar" is excluded.
	var got []string
	for key := range txn.IterKeys(kv.Key{"role"}, kv.KeyRange{
		From: kv.Key{"role", "bar"},
		To:   kv.Key{"role", "foo"},
	}) {
		got = append(got, key[1])
	}

	want := []string{"bar", "foo"}
	if !slices.Equal(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// TestIterPairsInclusiveTo mirrors TestIterKeysInclusiveTo for IterPairs.
func TestIterPairsInclusiveTo(t *testing.T) {
	txn := newTestTx(t)
	seedEntries(t, txn, "001", "002", "003", "004")

	var got []string
	for pair := range txn.IterPairs(kv.Key{"entry", "s"}, kv.KeyRange{
		From: kv.Key{"entry", "s", "002"},
		To:   kv.Key{"entry", "s", "003"},
	}) {
		got = append(got, pair.Key()[2])
	}

	want := []string{"002", "003"} // 003 (To) included, 004 excluded
	if !slices.Equal(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// TestKeyRoundTripWithColonSegment guards the lossless key encoding: a segment
// containing a ':' (e.g. a field named "com.acme:region") must round-trip
// through Set and iteration without being mis-split — which the old ':'
// separator would have done.
func TestKeyRoundTripWithColonSegment(t *testing.T) {
	txn := newTestTx(t)

	key := kv.Key{"stream", "field", "s", "com.acme:region"}
	if err := txn.Set(key, kv.Value("x")); err != nil {
		t.Fatalf("failed to set key: %v", err)
	}

	var got []kv.Key
	for pair := range txn.IterPairs(kv.Key{"stream", "field", "s"}, kv.KeyRange{}) {
		got = append(got, pair.Key())
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 key, got %d: %v", len(got), got)
	}
	if !slices.Equal(got[0], key) {
		t.Fatalf("expected %v, got %v", key, got[0])
	}
}
