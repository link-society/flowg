//go:build integration_fdb

package log_test

import (
	"testing"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/utils/langs/filtering"

	fdb_log "link-society.com/flowg/internal/storage/backends/foundationdb/concrete/log"
)

func connectString() string {
	return "docker:docker@127.0.0.1:4500"
}

func newLogStorage(t *testing.T) storage.LogStorage {
	t.Helper()

	opts := fdb_log.DefaultOptions()
	opts.ConnectionString = connectString()

	var logStorage storage.LogStorage

	app := fxtest.New(
		t,
		fdb_log.NewStorage(opts),
		fx.Populate(&logStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	t.Cleanup(app.RequireStop)

	return logStorage
}

func TestLog_All(t *testing.T) {
	ctx := t.Context()
	logStore := newLogStorage(t)

	// ========================================================================
	//  Stream Configs (empty)
	// ========================================================================
	t.Log("=== Stream Configs ===")

	streams, err := logStore.ListStreamConfigs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(streams) != 0 {
		t.Fatalf("expected 0 streams initially, got %d", len(streams))
	}

	// ========================================================================
	//  GetOrCreateStreamConfig
	// ========================================================================
	t.Log("=== GetOrCreateStreamConfig ===")

	cfg, err := logStore.GetOrCreateStreamConfig(ctx, "test-stream")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.RetentionTime != 0 {
		t.Fatalf("expected default RetentionTime 0, got %d", cfg.RetentionTime)
	}
	if cfg.RetentionSize != 0 {
		t.Fatalf("expected default RetentionSize 0, got %d", cfg.RetentionSize)
	}
	if len(cfg.IndexedFields) != 0 {
		t.Fatalf("expected empty IndexedFields, got %v", cfg.IndexedFields)
	}

	// Ensure second call returns same config (idempotent)
	cfg2, err := logStore.GetOrCreateStreamConfig(ctx, "test-stream")
	if err != nil {
		t.Fatal(err)
	}
	if cfg2.RetentionTime != cfg.RetentionTime {
		t.Fatal("GetOrCreateStreamConfig should be idempotent")
	}

	// List should now show 1 stream
	streams, err = logStore.ListStreamConfigs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(streams) != 1 {
		t.Fatalf("expected 1 stream, got %d", len(streams))
	}

	// ========================================================================
	//  ConfigureStream (indexing + retention)
	// ========================================================================
	t.Log("=== ConfigureStream ===")

	newCfg := models.StreamConfig{
		RetentionTime: 3600,       // 1 hour
		RetentionSize: 100,        // 100 MB
		IndexedFields: []string{"level", "app", "env"},
	}
	if err := logStore.ConfigureStream(ctx, "test-stream", newCfg); err != nil {
		t.Fatal(err)
	}

	cfg, err = logStore.GetOrCreateStreamConfig(ctx, "test-stream")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.RetentionTime != 3600 {
		t.Fatalf("expected RetentionTime 3600, got %d", cfg.RetentionTime)
	}
	if cfg.RetentionSize != 100 {
		t.Fatalf("expected RetentionSize 100, got %d", cfg.RetentionSize)
	}
	if len(cfg.IndexedFields) != 3 {
		t.Fatalf("expected 3 IndexedFields, got %d", len(cfg.IndexedFields))
	}

	// ========================================================================
	//  Ingest + FetchLogs
	// ========================================================================
	t.Log("=== Ingest + FetchLogs ===")

	now := time.Now().UTC()

	entry1 := &models.LogRecord{
		Timestamp: now.Add(-10 * time.Second),
		Fields: map[string]string{
			"level":   "info",
			"app":     "myapp",
			"env":     "prod",
			"message": "server started",
		},
	}
	key1, err := logStore.Ingest(ctx, "test-stream", entry1)
	if err != nil {
		t.Fatal(err)
	}
	if len(key1) == 0 {
		t.Fatal("expected non-empty key from Ingest")
	}

	entry2 := &models.LogRecord{
		Timestamp: now.Add(-5 * time.Second),
		Fields: map[string]string{
			"level":   "error",
			"app":     "myapp",
			"env":     "prod",
			"message": "connection refused",
		},
	}
	key2, err := logStore.Ingest(ctx, "test-stream", entry2)
	if err != nil {
		t.Fatal(err)
	}
	if len(key2) == 0 {
		t.Fatal("expected non-empty key from Ingest")
	}

	entry3 := &models.LogRecord{
		Timestamp: now.Add(-1 * time.Second),
		Fields: map[string]string{
			"level":   "info",
			"app":     "other-app",
			"env":     "staging",
			"message": "heartbeat",
		},
	}
	_, err = logStore.Ingest(ctx, "test-stream", entry3)
	if err != nil {
		t.Fatal(err)
	}

	// --- Fetch all logs in time range (no filter) ---
	from := now.Add(-30 * time.Second)
	to := now.Add(30 * time.Second)

	results, err := logStore.FetchLogs(ctx, "test-stream", from, to, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 log records, got %d", len(results))
	}

	// Most recent first (reverse order)
	if results[0].Fields["message"] != "heartbeat" {
		t.Fatalf("expected newest first, got: %s", results[0].Fields["message"])
	}

	// ========================================================================
	//  FetchLogs with indexing filter
	// ========================================================================
	t.Log("=== FetchLogs with index filter ===")

	// Filter by level=info (indexed field)
	results, err = logStore.FetchLogs(ctx, "test-stream", from, to, nil, map[string][]string{
		"level": {"info"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 info logs, got %d", len(results))
	}

	// Filter by level=error
	results, err = logStore.FetchLogs(ctx, "test-stream", from, to, nil, map[string][]string{
		"level": {"error"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 error log, got %d", len(results))
	}

	// Filter by level=warn (no matches)
	results, err = logStore.FetchLogs(ctx, "test-stream", from, to, nil, map[string][]string{
		"level": {"warn"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 warn logs, got %d", len(results))
	}

	// Filter by multiple values (OR within field): level=info OR level=error
	results, err = logStore.FetchLogs(ctx, "test-stream", from, to, nil, map[string][]string{
		"level": {"info", "error"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 logs for info+error, got %d", len(results))
	}

	// AND across fields: level=info AND app=myapp
	results, err = logStore.FetchLogs(ctx, "test-stream", from, to, nil, map[string][]string{
		"level": {"info"},
		"app":   {"myapp"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 log (info+myapp), got %d", len(results))
	}

	// ========================================================================
	//  FetchLogs with expr filter
	// ========================================================================
	t.Log("=== FetchLogs with expr filter ===")

	exprFilter, err := filtering.Compile(`level == "error"`)
	if err != nil {
		t.Fatal(err)
	}

	results, err = logStore.FetchLogs(ctx, "test-stream", from, to, exprFilter, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 error log from expr filter, got %d", len(results))
	}
	if results[0].Fields["message"] != "connection refused" {
		t.Fatalf("unexpected log: %v", results[0].Fields)
	}

	// Combined: index filter + expr filter
	results, err = logStore.FetchLogs(ctx, "test-stream", from, to, exprFilter, map[string][]string{
		"app": {"myapp"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 error+myapp log, got %d", len(results))
	}

	// ========================================================================
	//  Empty time range
	// ========================================================================
	t.Log("=== Empty time range ===")

	futureFrom := now.Add(24 * time.Hour)
	futureTo := now.Add(48 * time.Hour)
	results, err = logStore.FetchLogs(ctx, "test-stream", futureFrom, futureTo, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 logs in future range, got %d", len(results))
	}

	// ========================================================================
	//  ListStreamFields
	// ========================================================================
	t.Log("=== ListStreamFields ===")

	fields, err := logStore.ListStreamFields(ctx, "test-stream")
	if err != nil {
		t.Fatal(err)
	}
	fieldSet := make(map[string]bool)
	for _, f := range fields {
		fieldSet[f] = true
	}
	for _, expected := range []string{"level", "app", "env", "message"} {
		if !fieldSet[expected] {
			t.Fatalf("ListStreamFields missing expected field: %s", expected)
		}
	}

	// ========================================================================
	//  Distinct (indexed field values)
	// ========================================================================
	t.Log("=== Distinct ===")

	distinct, err := logStore.Distinct(ctx, "test-stream")
	if err != nil {
		t.Fatal(err)
	}

	// level should have: info, error
	if len(distinct["level"]) != 2 {
		t.Fatalf("expected 2 distinct levels, got %d: %v", len(distinct["level"]), distinct["level"])
	}
	// app should have: myapp, other-app
	if len(distinct["app"]) != 2 {
		t.Fatalf("expected 2 distinct apps, got %d: %v", len(distinct["app"]), distinct["app"])
	}
	// env should have: prod, staging
	if len(distinct["env"]) != 2 {
		t.Fatalf("expected 2 distinct envs, got %d: %v", len(distinct["env"]), distinct["env"])
	}

	// ========================================================================
	//  StreamUsage
	// ========================================================================
	t.Log("=== StreamUsage ===")

	usage, err := logStore.StreamUsage(ctx, "test-stream")
	if err != nil {
		t.Fatal(err)
	}
	// FDB GetEstimatedRangeSizeBytes may return 0 for very small data;
	// any non-negative value is valid — the important thing is no error.
	if usage < 0 {
		t.Fatalf("expected non-negative StreamUsage, got %d", usage)
	}

	// Non-existent stream
	usage, err = logStore.StreamUsage(ctx, "no-such-stream")
	if err != nil {
		t.Fatal(err)
	}
	if usage != 0 {
		t.Fatalf("expected 0 usage for non-existent stream, got %d", usage)
	}

	// ========================================================================
	//  IndexField / UnindexField
	// ========================================================================
	t.Log("=== IndexField / UnindexField ===")

	// Add a new indexed field after ingestion
	if err := logStore.IndexField(ctx, "test-stream", "message"); err != nil {
		t.Fatal(err)
	}

	// Distinct should now include message
	distinct, err = logStore.Distinct(ctx, "test-stream")
	if err != nil {
		t.Fatal(err)
	}
	if len(distinct["message"]) != 3 {
		t.Fatalf("expected 3 distinct messages after IndexField, got %d: %v", len(distinct["message"]), distinct["message"])
	}

	// Remove an indexed field
	if err := logStore.UnindexField(ctx, "test-stream", "message"); err != nil {
		t.Fatal(err)
	}

	// Distinct should no longer include message
	distinct, err = logStore.Distinct(ctx, "test-stream")
	if err != nil {
		t.Fatal(err)
	}
	if _, exists := distinct["message"]; exists {
		t.Fatal("expected message to be removed from distinct after UnindexField")
	}

	// ========================================================================
	//  DeleteStream (cascades entries, fields, config, index)
	// ========================================================================
	t.Log("=== DeleteStream ===")

	if err := logStore.DeleteStream(ctx, "test-stream"); err != nil {
		t.Fatal(err)
	}

	// Stream should no longer appear in list
	streams, err = logStore.ListStreamConfigs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(streams) != 0 {
		t.Fatalf("expected 0 streams after delete, got %d", len(streams))
	}

	// Fetch logs should return empty
	results, err = logStore.FetchLogs(ctx, "test-stream", from, to, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 logs after stream deletion, got %d", len(results))
	}

	// Fields should be gone
	fields, err = logStore.ListStreamFields(ctx, "test-stream")
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 0 {
		t.Fatalf("expected 0 fields after stream deletion, got %d", len(fields))
	}

	// ========================================================================
	//  Multiple streams (isolation)
	// ========================================================================
	t.Log("=== Multi-stream isolation ===")

	logStore.GetOrCreateStreamConfig(ctx, "stream-a")
	logStore.GetOrCreateStreamConfig(ctx, "stream-b")

	logStore.Ingest(ctx, "stream-a", &models.LogRecord{
		Timestamp: now,
		Fields:    map[string]string{"app": "a-only"},
	})
	logStore.Ingest(ctx, "stream-b", &models.LogRecord{
		Timestamp: now,
		Fields:    map[string]string{"app": "b-only"},
	})

	aResults, err := logStore.FetchLogs(ctx, "stream-a", from, to, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(aResults) != 1 || aResults[0].Fields["app"] != "a-only" {
		t.Fatal("stream isolation failed: stream-a has wrong data")
	}

	bResults, err := logStore.FetchLogs(ctx, "stream-b", from, to, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(bResults) != 1 || bResults[0].Fields["app"] != "b-only" {
		t.Fatal("stream isolation failed: stream-b has wrong data")
	}

	// ========================================================================
	//  ConfigureStream changes indexing (back-fill + drop)
	// ========================================================================
	t.Log("=== ConfigureStream index back-fill ===")

	// Configure stream-b with indexing on "app"
	logStore.ConfigureStream(ctx, "stream-b", models.StreamConfig{
		RetentionTime: 0,
		RetentionSize: 0,
		IndexedFields: []string{"app"},
	})

	// Add more data
	logStore.Ingest(ctx, "stream-b", &models.LogRecord{
		Timestamp: now,
		Fields:    map[string]string{"app": "b-only", "extra": "true"},
	})

	// Should be findable via index
	filteredB, err := logStore.FetchLogs(ctx, "stream-b", from, to, nil, map[string][]string{
		"app": {"b-only"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(filteredB) != 2 {
		t.Fatalf("expected 2 logs in stream-b via index filter, got %d", len(filteredB))
	}

	// Remove app from indexed fields via ConfigureStream
	logStore.ConfigureStream(ctx, "stream-b", models.StreamConfig{
		IndexedFields: []string{},
	})

	// Distinct should no longer have app
	distinctB, err := logStore.Distinct(ctx, "stream-b")
	if err != nil {
		t.Fatal(err)
	}
	if _, exists := distinctB["app"]; exists {
		t.Fatal("expected app index to be removed after ConfigureStream removed it from IndexedFields")
	}
}
