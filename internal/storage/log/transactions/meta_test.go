package transactions_test

import (
	"testing"

	"github.com/dgraph-io/badger/v4"
	badgerOptions "github.com/dgraph-io/badger/v4/options"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/log/transactions"
	"link-society.com/flowg/internal/utils/hlc"
)

func newDB(t *testing.T) *badger.DB {
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

// TestConfigureStreamPersistsOnFreshStream guards against a regression where
// configuring a brand new stream silently dropped the configuration.
//
// ConfigureStream used to call GetOrCreateStreamConfig, which auto-creates an
// empty LWW envelope stamped at `ts`. The subsequent LWW Apply of the real
// config was also stamped at `ts`, and because LWW only accepts strictly newer
// timestamps, the real config was rejected and the stream kept a zero-value
// configuration (e.g. RetentionTime = 0, meaning logs never expired).
func TestConfigureStreamPersistsOnFreshStream(t *testing.T) {
	db := newDB(t)

	ts := hlc.Timestamp{WallTime: 1000, Logical: 0, NodeID: "node-a"}
	want := models.StreamConfig{
		RetentionTime: 2,
		RetentionSize: 5,
		IndexedFields: []string{"level"},
	}

	err := db.Update(func(txn *badger.Txn) error {
		return transactions.ConfigureStream(txn, "mystream", want, ts)
	})
	if err != nil {
		t.Fatalf("ConfigureStream: %v", err)
	}

	var got models.StreamConfig
	readTs := hlc.Timestamp{WallTime: 2000, Logical: 0, NodeID: "node-a"}
	err = db.Update(func(txn *badger.Txn) error {
		cfg, _, err := transactions.GetOrCreateStreamConfig(txn, "mystream", readTs)
		got = cfg
		return err
	})
	if err != nil {
		t.Fatalf("GetOrCreateStreamConfig: %v", err)
	}

	if got.RetentionTime != want.RetentionTime {
		t.Errorf("RetentionTime = %d; want %d", got.RetentionTime, want.RetentionTime)
	}
	if got.RetentionSize != want.RetentionSize {
		t.Errorf("RetentionSize = %d; want %d", got.RetentionSize, want.RetentionSize)
	}
	if !got.IsFieldIndexed("level") {
		t.Errorf("field 'level' should be indexed; got IndexedFields=%v", got.IndexedFields)
	}
}
