package badger

import (
	"testing"

	"strings"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/storage/generic/kv"
)

func newTestDB(t *testing.T) *badger.DB {
	t.Helper()

	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true).WithLoggingLevel(badger.ERROR),
	)
	if err != nil {
		t.Fatalf("failed to open in-memory badger: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return db
}

// TestMigrateKeySeparator checks that legacy ':'-joined keys are re-encoded with
// the ESC separator (values preserved), and that the migration is idempotent and
// leaves post-migration keys that legitimately contain a ':' untouched.
func TestMigrateKeySeparator(t *testing.T) {
	db := newTestDB(t)

	// Seed legacy ':'-joined keys directly, as an older FlowG would have stored
	// them (including a field value that itself contained a ':').
	legacy := map[string]string{
		"role:admin:read_streams":        "",
		"stream:field:s:com.acme:region": "marker",
		"entry:s:0001:uuid":              `{"m":1}`,
	}
	err := db.Update(func(txn *badger.Txn) error {
		for k, v := range legacy {
			if err := txn.Set([]byte(k), []byte(v)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to seed legacy keys: %v", err)
	}

	if err := migrateKeySeparator(db); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	// Each legacy key must be gone, and its ESC-joined equivalent present with
	// the same value.
	err = db.View(func(txn *badger.Txn) error {
		for k, want := range legacy {
			if _, err := txn.Get([]byte(k)); err != badger.ErrKeyNotFound {
				t.Fatalf("legacy key %q still present (err=%v)", k, err)
			}

			newKey := []byte(strings.ReplaceAll(k, legacyKeySeparator, keySeparator))
			item, err := txn.Get(newKey)
			if err != nil {
				t.Fatalf("migrated key for %q missing: %v", k, err)
			}
			got, err := item.ValueCopy(nil)
			if err != nil {
				t.Fatalf("failed to read migrated value for %q: %v", k, err)
			}
			if string(got) != want {
				t.Fatalf("value mismatch for %q: got %q, want %q", k, got, want)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("verification failed: %v", err)
	}

	// A key written after migration may legitimately hold a ':' inside a segment.
	// Running the migration again must be a no-op (guarded by the marker) and
	// must not corrupt it.
	liveKey := kv.Key{"stream", "field", "s", "com.acme:region"}
	err = db.Update(func(txn *badger.Txn) error {
		return (&BadgerTx{concrete: txn}).Set(liveKey, kv.Value("live"))
	})
	if err != nil {
		t.Fatalf("failed to write live key: %v", err)
	}

	if err := migrateKeySeparator(db); err != nil {
		t.Fatalf("second migration failed: %v", err)
	}

	err = db.View(func(txn *badger.Txn) error {
		got, err := (&BadgerTx{concrete: txn}).Get(liveKey)
		if err != nil {
			return err
		}
		if string(got) != "live" {
			t.Fatalf("live key was corrupted by re-migration: got %q", got)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("live key verification failed: %v", err)
	}
}
