//go:build integration_fdb

package kvstore_test

import (
	"bytes"
	"strings"
	"testing"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	fdbkvstore "link-society.com/flowg/internal/storage/backends/foundationdb/kvstore"
)

func connectString() string {
	return "docker:docker@127.0.0.1:4500"
}

func TestKVStore_All(t *testing.T) {
	ctx := t.Context()

	opts := fdbkvstore.DefaultOptions()
	opts.ConnectionString = connectString()
	opts.Prefix = []byte("kvtest/")

	type deps struct {
		fx.In
		Store fdbkvstore.Storage `name:"kv"`
	}

	var d deps
	app := fxtest.New(t,
		fdbkvstore.NewStorage(opts),
		fx.Populate(&d),
		fx.NopLogger,
	)
	store := d.Store
	app.RequireStart()
	defer app.RequireStop()

	// --- View on empty store ---
	err := store.View(ctx, func(tr fdb.ReadTransaction) error {
		ri := tr.GetRange(fdb.KeyRange{
			Begin: fdb.Key("kvtest/"),
			End:   fdb.Key("kvtest0"),
		}, fdb.RangeOptions{}).Iterator()
		if ri.Advance() {
			t.Fatal("expected no keys initially")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// --- Write data via Update ---
	err = store.Update(ctx, func(tr fdb.Transaction) error {
		tr.Set(fdb.Key("kvtest/a"), []byte("value-a"))
		tr.Set(fdb.Key("kvtest/b"), []byte("value-b"))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// --- Read via View ---
	err = store.View(ctx, func(tr fdb.ReadTransaction) error {
		val, err := tr.Get(fdb.Key("kvtest/a")).Get()
		if err != nil {
			t.Fatal(err)
		}
		if string(val) != "value-a" {
			t.Fatalf("expected value-a, got %s", val)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// --- Backup ---
	var buf bytes.Buffer
	ver, err := store.Backup(ctx, &buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if ver == 0 {
		t.Fatal("expected non-zero version")
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty backup")
	}

	backupContent := buf.String()
	if len(backupContent) > 200 {
		t.Logf("Backup content (first 200 chars): %s", backupContent[:200])
	} else {
		t.Logf("Backup content: %s", backupContent)
	}
	// The backup encodes values as base64
	if !strings.Contains(backupContent, "dmFsdWUtYQ==") {
		t.Fatal("backup should contain our data (value-a base64)")
	}

	// --- Incremental backup (since=current ver) should be empty or minimal ---
	buf.Reset()
	_, err = store.Backup(ctx, &buf, ver)
	if err != nil {
		t.Fatal(err)
	}
	// May or may not contain data depending on FDB read version stability

	// --- Clear the prefix ---
	err = store.Update(ctx, func(tr fdb.Transaction) error {
		tr.ClearRange(fdb.KeyRange{
			Begin: fdb.Key("kvtest/"),
			End:   fdb.Key("kvtest0"),
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify empty
	err = store.View(ctx, func(tr fdb.ReadTransaction) error {
		val, err := tr.Get(fdb.Key("kvtest/a")).Get()
		if err != nil {
			t.Fatal(err)
		}
		if val != nil {
			t.Fatal("expected nil after clear")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// --- Restore from original backup ---
	err = store.Restore(ctx, strings.NewReader(backupContent))
	if err != nil {
		t.Fatal(err)
	}

	// Verify data back
	err = store.View(ctx, func(tr fdb.ReadTransaction) error {
		val, err := tr.Get(fdb.Key("kvtest/a")).Get()
		if err != nil {
			t.Fatal(err)
		}
		if string(val) != "value-a" {
			t.Fatalf("expected value-a after restore, got %s", val)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
