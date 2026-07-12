package badger

import (
	"errors"
	"strings"
	"testing"
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
