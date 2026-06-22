package schema

import (
	"testing"

	"github.com/dgraph-io/badger/v4"
	badgerOptions "github.com/dgraph-io/badger/v4/options"

	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/lww"
)

func newDB(t *testing.T) *badger.DB {
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
		t.Fatalf("open badger: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func migrationTimestamp() hlc.Timestamp {
	return hlc.Timestamp{WallTime: 42, Logical: 0, NodeID: "test"}
}

func TestReadVersionAbsent(t *testing.T) {
	t.Parallel()
	db := newDB(t)

	err := db.View(func(txn *badger.Txn) error {
		version, err := readVersion(txn)
		if err != nil {
			return err
		}
		if version != 0 {
			t.Fatalf("expected version 0, got %d", version)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("view: %v", err)
	}
}

func TestWriteReadVersionRoundTrip(t *testing.T) {
	t.Parallel()
	db := newDB(t)

	err := db.Update(func(txn *badger.Txn) error {
		return writeVersion(txn, 1)
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	err = db.View(func(txn *badger.Txn) error {
		version, err := readVersion(txn)
		if err != nil {
			return err
		}
		if version != 1 {
			t.Fatalf("expected version 1, got %d", version)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("view: %v", err)
	}
}

func TestEnvelopeV0toV1WrapsAllKeys(t *testing.T) {
	t.Parallel()
	db := newDB(t)

	raw := map[string]string{
		"role:admin:read_acls": "",
		"user:root:password":   "hashed-password",
		"index:user:root":      "",
	}

	err := db.Update(func(txn *badger.Txn) error {
		for k, v := range raw {
			if err := txn.Set([]byte(k), []byte(v)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		return envelopeV0toV1(txn, migrationTimestamp(), nil)
	})
	if err != nil {
		t.Fatalf("envelope: %v", err)
	}

	err = db.View(func(txn *badger.Txn) error {
		for k, want := range raw {
			env, found, err := lww.Read(txn, []byte(k))
			if err != nil {
				return err
			}
			if !found {
				t.Fatalf("key %q not found after migration", k)
			}
			if string(env.Payload) != want {
				t.Fatalf("key %q: expected payload %q, got %q", k, want, env.Payload)
			}
			if env.Timestamp != (hlc.Timestamp{}) {
				t.Fatalf("key %q: expected zero sentinel timestamp, got %v", k, env.Timestamp)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestEnvelopeV0toV1SkipsVersionKey(t *testing.T) {
	t.Parallel()
	db := newDB(t)

	err := db.Update(func(txn *badger.Txn) error {
		if err := writeVersion(txn, 0); err != nil {
			return err
		}
		return txn.Set([]byte("role:admin:read_acls"), []byte(""))
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		return envelopeV0toV1(txn, migrationTimestamp(), nil)
	})
	if err != nil {
		t.Fatalf("envelope: %v", err)
	}

	err = db.View(func(txn *badger.Txn) error {
		version, err := readVersion(txn)
		if err != nil {
			return err
		}
		if version != 0 {
			t.Fatalf("version key should be untouched (raw 0), got %d", version)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestEnvelopeV0toV1RespectsPrefixes(t *testing.T) {
	t.Parallel()
	db := newDB(t)

	err := db.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte("stream:config:default"), []byte(`{"retention":1}`)); err != nil {
			return err
		}
		return txn.Set([]byte("entry:default:00000001"), []byte("raw-log-entry"))
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		return envelopeV0toV1(txn, migrationTimestamp(), [][]byte{[]byte("stream:config:")})
	})
	if err != nil {
		t.Fatalf("envelope: %v", err)
	}

	err = db.View(func(txn *badger.Txn) error {
		env, found, err := lww.Read(txn, []byte("stream:config:default"))
		if err != nil {
			return err
		}
		if !found || string(env.Payload) != `{"retention":1}` {
			t.Fatalf("stream config not enveloped correctly: found=%v payload=%q", found, env.Payload)
		}

		item, err := txn.Get([]byte("entry:default:00000001"))
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		if string(val) != "raw-log-entry" {
			t.Fatalf("log entry should be untouched, got %q", val)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
}
