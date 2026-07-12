package log_test

import (
	"testing"

	"strings"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/cmd/flowg-server/logging"
	"link-society.com/flowg/internal/models"

	badgerlog "link-society.com/flowg/internal/storage/backends/badger/concrete/log"
	"link-society.com/flowg/internal/storage/generic/kv"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// TestFetchLogsIndexingFilter guards against the indexed-field filter matching
// every record (which happens when the "missing key" case is not distinguished
// from the "present, empty value" case on the kv.QueryTx.Get contract).
func TestFetchLogsIndexingFilter(t *testing.T) {
	logging.Discard()

	ctx := t.Context()

	opts := badgerlog.DefaultOptions()
	opts.InMemory = true

	var logStorage storage.LogStorage

	app := fxtest.New(
		t,
		badgerlog.NewStorage(opts),
		fx.Populate(&logStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	const stream = "test"

	if err := logStorage.ConfigureStream(ctx, stream, models.StreamConfig{
		IndexedFields: []string{"level"},
	}); err != nil {
		t.Fatalf("failed to configure stream: %v", err)
	}

	records := []map[string]string{
		{"level": "error", "message": "boom"},
		{"level": "info", "message": "hello"},
		{"level": "info", "message": "world"},
	}
	for _, fields := range records {
		if _, err := logStorage.Ingest(ctx, stream, models.NewLogRecord(fields)); err != nil {
			t.Fatalf("failed to ingest record: %v", err)
		}
	}

	from := time.Now().Add(-time.Minute)
	to := time.Now().Add(time.Minute)

	// Sanity: without indexing, every record in the window is returned.
	all, err := logStorage.FetchLogs(ctx, stream, from, to, nil, nil)
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 records without indexing, got %d", len(all))
	}

	// Filtering on an indexed value must return ONLY the matching record.
	errorOnly, err := logStorage.FetchLogs(ctx, stream, from, to, nil, map[string][]string{
		"level": {"error"},
	})
	if err != nil {
		t.Fatalf("failed to fetch logs with indexing: %v", err)
	}
	if len(errorOnly) != 1 {
		t.Fatalf("indexing {level:[error]}: expected 1 record, got %d", len(errorOnly))
	}
	if errorOnly[0].Fields["level"] != "error" {
		t.Fatalf("expected the error record, got level=%q", errorOnly[0].Fields["level"])
	}

	// A value that no record carries must match nothing.
	none, err := logStorage.FetchLogs(ctx, stream, from, to, nil, map[string][]string{
		"level": {"debug"},
	})
	if err != nil {
		t.Fatalf("failed to fetch logs with indexing: %v", err)
	}
	if len(none) != 0 {
		t.Fatalf("indexing {level:[debug]}: expected 0 records, got %d", len(none))
	}
}

// TestStreamIndicesDistinct guards against Distinct() skipping every real index
// key (which happens when the segment-count guard is wrong), returning an empty
// map and breaking the stream-indices autocomplete endpoint.
func TestStreamIndicesDistinct(t *testing.T) {
	logging.Discard()

	ctx := t.Context()

	opts := badgerlog.DefaultOptions()
	opts.InMemory = true

	var logStorage storage.LogStorage

	app := fxtest.New(
		t,
		badgerlog.NewStorage(opts),
		fx.Populate(&logStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	const stream = "test"

	if err := logStorage.ConfigureStream(ctx, stream, models.StreamConfig{
		IndexedFields: []string{"level"},
	}); err != nil {
		t.Fatalf("failed to configure stream: %v", err)
	}

	for _, fields := range []map[string]string{
		{"level": "error", "message": "boom"},
		{"level": "info", "message": "hello"},
		{"level": "info", "message": "world"},
	} {
		if _, err := logStorage.Ingest(ctx, stream, models.NewLogRecord(fields)); err != nil {
			t.Fatalf("failed to ingest record: %v", err)
		}
	}

	indices, err := logStorage.Distinct(ctx, stream)
	if err != nil {
		t.Fatalf("failed to get stream indices: %v", err)
	}

	values := indices["level"]
	if len(values) != 2 {
		t.Fatalf("expected 2 distinct values for 'level', got %d (%v)", len(values), values)
	}

	seen := map[string]bool{}
	for _, v := range values {
		seen[v] = true
	}
	if !seen["error"] || !seen["info"] {
		t.Fatalf("expected distinct 'level' values to include error and info, got %v", values)
	}
}

// TestIngestOversizedIndexedValue guards against an indexed field whose value is
// too large to fit in an index key aborting the whole ingestion. The record must
// still be stored and queryable by time; only the oversized value is left
// unindexed, so an exact-match filter on it degrades to "no match" instead of
// erroring.
func TestIngestOversizedIndexedValue(t *testing.T) {
	logging.Discard()

	ctx := t.Context()

	opts := badgerlog.DefaultOptions()
	opts.InMemory = true

	var logStorage storage.LogStorage

	app := fxtest.New(
		t,
		badgerlog.NewStorage(opts),
		fx.Populate(&logStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	const stream = "test"

	if err := logStorage.ConfigureStream(ctx, stream, models.StreamConfig{
		IndexedFields: []string{"body"},
	}); err != nil {
		t.Fatalf("failed to configure stream: %v", err)
	}

	// A value large enough that its base64-encoded form overflows the key-size
	// limit once embedded in the index key.
	hugeValue := strings.Repeat("x", kv.MaxKeySize)

	// Ingestion must succeed despite the oversized indexed value.
	if _, err := logStorage.Ingest(ctx, stream, models.NewLogRecord(map[string]string{
		"body": hugeValue,
	})); err != nil {
		t.Fatalf("ingestion of oversized indexed value should not fail: %v", err)
	}

	from := time.Now().Add(-time.Minute)
	to := time.Now().Add(time.Minute)

	// The record is still stored and queryable by time.
	all, err := logStorage.FetchLogs(ctx, stream, from, to, nil, nil)
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected the record to be stored, got %d records", len(all))
	}

	// Filtering by the oversized value must degrade to no match, not error.
	filtered, err := logStorage.FetchLogs(ctx, stream, from, to, nil, map[string][]string{
		"body": {hugeValue},
	})
	if err != nil {
		t.Fatalf("filtering by oversized value should not error: %v", err)
	}
	if len(filtered) != 0 {
		t.Fatalf("expected 0 records for oversized filter, got %d", len(filtered))
	}
}
