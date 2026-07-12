package log_test

import (
	"testing"

	"context"
	"fmt"

	"strings"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/cmd/flowg-server/logging"
	"link-society.com/flowg/internal/models"

	badgerlog "link-society.com/flowg/internal/storage/backends/badger/concrete/log"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// newBatchedStorage builds an in-memory log storage whose whole-stream
// operations batch at batchSize keys per transaction, so tests can force the
// multi-transaction paths with only a handful of records.
func newBatchedStorage(t *testing.T, batchSize int) (context.Context, storage.LogStorage) {
	t.Helper()
	logging.Discard()

	opts := badgerlog.DefaultOptions()
	opts.InMemory = true
	opts.BatchSize = batchSize

	var logStorage storage.LogStorage

	app := fxtest.New(
		t,
		badgerlog.NewStorage(opts),
		fx.Populate(&logStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	t.Cleanup(app.RequireStop)

	return t.Context(), logStorage
}

// TestDeleteStreamBatched guards against DeleteStream aborting on a large stream
// (which happens when it tries to clear every key in one transaction): with a
// tiny batch size the deletion must span several transactions and still remove
// every entry, index key, field marker and the config.
func TestDeleteStreamBatched(t *testing.T) {
	ctx, logStorage := newBatchedStorage(t, 3)

	const stream = "test"

	if err := logStorage.ConfigureStream(ctx, stream, models.StreamConfig{
		IndexedFields: []string{"level"},
	}); err != nil {
		t.Fatalf("failed to configure stream: %v", err)
	}

	for i := 0; i < 7; i++ {
		if _, err := logStorage.Ingest(ctx, stream, models.NewLogRecord(map[string]string{
			"level": fmt.Sprintf("l%d", i),
		})); err != nil {
			t.Fatalf("failed to ingest record: %v", err)
		}
	}

	// Sanity: the stream exists before deletion.
	configs, err := logStorage.ListStreamConfigs(ctx)
	if err != nil {
		t.Fatalf("failed to list stream configs: %v", err)
	}
	if _, ok := configs[stream]; !ok {
		t.Fatalf("expected stream %q to exist before deletion", stream)
	}

	if err := logStorage.DeleteStream(ctx, stream); err != nil {
		t.Fatalf("failed to delete stream: %v", err)
	}

	configs, err = logStorage.ListStreamConfigs(ctx)
	if err != nil {
		t.Fatalf("failed to list stream configs: %v", err)
	}
	if _, ok := configs[stream]; ok {
		t.Fatalf("stream config was not deleted")
	}

	fields, err := logStorage.ListStreamFields(ctx, stream)
	if err != nil {
		t.Fatalf("failed to list stream fields: %v", err)
	}
	if len(fields) != 0 {
		t.Fatalf("expected no fields after deletion, got %v", fields)
	}

	indices, err := logStorage.Distinct(ctx, stream)
	if err != nil {
		t.Fatalf("failed to get stream indices: %v", err)
	}
	if len(indices) != 0 {
		t.Fatalf("expected no indices after deletion, got %v", indices)
	}

	from := time.Now().Add(-time.Minute)
	to := time.Now().Add(time.Minute)
	recs, err := logStorage.FetchLogs(ctx, stream, from, to, nil, nil)
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	if len(recs) != 0 {
		t.Fatalf("expected no records after deletion, got %d", len(recs))
	}
}

// TestConfigureStreamBackfillAndUnindexBatched guards against the index back-
// fill (and its removal) aborting on a large stream: with a tiny batch size both
// must span several transactions and still index / unindex every record.
func TestConfigureStreamBackfillAndUnindexBatched(t *testing.T) {
	ctx, logStorage := newBatchedStorage(t, 3)

	const stream = "test"

	// Ingest before indexing, so the index must be back-filled afterwards.
	levels := []string{"error", "info", "info", "error", "info", "error", "info"}
	for _, level := range levels {
		if _, err := logStorage.Ingest(ctx, stream, models.NewLogRecord(map[string]string{
			"level": level,
		})); err != nil {
			t.Fatalf("failed to ingest record: %v", err)
		}
	}

	// Enabling the index triggers a batched back-fill over the 7 existing records.
	if err := logStorage.ConfigureStream(ctx, stream, models.StreamConfig{
		IndexedFields: []string{"level"},
	}); err != nil {
		t.Fatalf("failed to configure stream: %v", err)
	}

	from := time.Now().Add(-time.Minute)
	to := time.Now().Add(time.Minute)

	errRecs, err := logStorage.FetchLogs(ctx, stream, from, to, nil, map[string][]string{
		"level": {"error"},
	})
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	if len(errRecs) != 3 {
		t.Fatalf("expected 3 back-filled error records, got %d", len(errRecs))
	}

	infoRecs, err := logStorage.FetchLogs(ctx, stream, from, to, nil, map[string][]string{
		"level": {"info"},
	})
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	if len(infoRecs) != 4 {
		t.Fatalf("expected 4 back-filled info records, got %d", len(infoRecs))
	}

	// Disabling the index triggers a batched removal of every index key.
	if err := logStorage.ConfigureStream(ctx, stream, models.StreamConfig{
		IndexedFields: []string{},
	}); err != nil {
		t.Fatalf("failed to reconfigure stream: %v", err)
	}

	indices, err := logStorage.Distinct(ctx, stream)
	if err != nil {
		t.Fatalf("failed to get stream indices: %v", err)
	}
	if len(indices) != 0 {
		t.Fatalf("expected no indices after unindexing, got %v", indices)
	}

	// The records themselves are untouched by (un)indexing.
	all, err := logStorage.FetchLogs(ctx, stream, from, to, nil, nil)
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	if len(all) != len(levels) {
		t.Fatalf("expected %d records to remain, got %d", len(levels), len(all))
	}
}

// TestCollectGarbageBatched guards against retention enforcement aborting (or
// silently stopping) on a large stream: with a tiny batch size the eviction must
// span several transactions, bring the stream within budget, and purge the index
// references of every evicted record in lockstep.
func TestCollectGarbageBatched(t *testing.T) {
	ctx, logStorage := newBatchedStorage(t, 3)

	collector, ok := logStorage.(interface {
		CollectGarbage(context.Context) error
	})
	if !ok {
		t.Fatalf("log storage does not expose CollectGarbage")
	}

	const stream = "test"
	const retentionMB = int64(1)

	if err := logStorage.ConfigureStream(ctx, stream, models.StreamConfig{
		IndexedFields: []string{"seq"},
		RetentionSize: retentionMB,
	}); err != nil {
		t.Fatalf("failed to configure stream: %v", err)
	}

	// 40 x ~50 KB comfortably exceeds the 1 MB budget.
	blob := strings.Repeat("x", 50_000)
	const n = 40
	for i := 0; i < n; i++ {
		if _, err := logStorage.Ingest(ctx, stream, models.NewLogRecord(map[string]string{
			"seq":  fmt.Sprintf("%03d", i),
			"blob": blob,
		})); err != nil {
			t.Fatalf("failed to ingest record: %v", err)
		}
	}

	budget := retentionMB * 1024 * 1024

	usage, err := logStorage.StreamUsage(ctx, stream)
	if err != nil {
		t.Fatalf("failed to estimate usage: %v", err)
	}
	if usage <= budget {
		t.Fatalf("test setup: stream is not over budget (%d <= %d)", usage, budget)
	}

	if err := collector.CollectGarbage(ctx); err != nil {
		t.Fatalf("failed to collect garbage: %v", err)
	}

	usage, err = logStorage.StreamUsage(ctx, stream)
	if err != nil {
		t.Fatalf("failed to estimate usage: %v", err)
	}
	if usage > budget {
		t.Fatalf("expected stream within budget after GC, got %d > %d", usage, budget)
	}

	from := time.Now().Add(-time.Minute)
	to := time.Now().Add(time.Minute)
	recs, err := logStorage.FetchLogs(ctx, stream, from, to, nil, nil)
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	if len(recs) == 0 || len(recs) == n {
		t.Fatalf("expected partial eviction, got %d of %d records", len(recs), n)
	}

	// Every evicted record's index reference must have been purged, and every
	// survivor's kept: the distinct 'seq' values must match the survivors exactly.
	indices, err := logStorage.Distinct(ctx, stream)
	if err != nil {
		t.Fatalf("failed to get stream indices: %v", err)
	}
	if len(indices["seq"]) != len(recs) {
		t.Fatalf(
			"index/entry mismatch after GC: %d distinct 'seq' values vs %d records",
			len(indices["seq"]), len(recs),
		)
	}
}
