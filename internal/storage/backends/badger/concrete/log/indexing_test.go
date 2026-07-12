package log_test

import (
	"testing"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/cmd/flowg-server/logging"
	"link-society.com/flowg/internal/models"

	badgerlog "link-society.com/flowg/internal/storage/backends/badger/concrete/log"
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
