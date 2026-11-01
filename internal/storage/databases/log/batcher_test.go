package log

import (
	"testing"

	"context"
	"errors"
	"log/slog"

	"encoding/json"
	"strings"
	"sync/atomic"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"go.link-society.com/flowg/pkg/models"

	"go.link-society.com/flowg/internal/storage/backends/badger"
	"go.link-society.com/flowg/internal/storage/generic/kv"
)

// countingAdapter counts the transactions opened through Update.
type countingAdapter struct {
	*badger.BadgerAdapter
	updates atomic.Int32
}

func (a *countingAdapter) Update(ctx context.Context, fn func(txn *badger.BadgerTx) error) error {
	a.updates.Add(1)
	return a.BadgerAdapter.Update(ctx, fn)
}

type testAdapterOptions struct {
	memTableSize int64
}

// newTestStorage builds a log storage over an in-memory BadgerDB whose memtable
// size bounds the transaction size (about 15% of it).
func newTestStorage(t *testing.T, memTableSize int64) (*Storage[*badger.BadgerTx, *badger.BadgerTx], *countingAdapter) {
	t.Helper()

	var concrete *badger.BadgerAdapter

	app := fxtest.New(
		t,
		fx.Supply(slog.New(slog.DiscardHandler)),
		fx.Supply(testAdapterOptions{memTableSize: memTableSize}),
		badger.NewAdapter("test", func(o testAdapterOptions) badger.AdapterOptions {
			return badger.AdapterOptions{InMemory: true, MemTableSize: o.memTableSize}
		}),
		fx.Invoke(fx.Annotate(
			func(a *badger.BadgerAdapter) { concrete = a },
			fx.ParamTags(`name:"test"`),
		)),
		fx.NopLogger,
	)
	app.RequireStart()
	t.Cleanup(app.RequireStop)

	adapter := &countingAdapter{BadgerAdapter: concrete}
	return NewStorage(slog.New(slog.DiscardHandler), adapter, 0, nil), adapter
}

// newTestRequest builds an Ingest call for a record holding fields, with the
// given encoded value (the JSON encoding of the record when val is nil).
func newTestRequest(t *testing.T, stream string, fields map[string]string, val []byte) (ingestRequest, <-chan error) {
	t.Helper()

	record := models.NewLogRecord(fields)
	if val == nil {
		var err error
		val, err = json.Marshal(record)
		if err != nil {
			t.Fatalf("failed to marshal record: %v", err)
		}
	}

	replyTo := make(chan error, 1)
	return ingestRequest{
		ctx:     t.Context(),
		stream:  stream,
		record:  record,
		key:     newDbKey(stream, record),
		val:     val,
		replyTo: replyTo,
	}, replyTo
}

func countRecords(t *testing.T, s *Storage[*badger.BadgerTx, *badger.BadgerTx], stream string) int {
	t.Helper()

	records, err := s.FetchLogs(t.Context(), stream, time.Now().Add(-time.Hour), time.Now().Add(time.Hour), nil, nil)
	if err != nil {
		t.Fatalf("failed to fetch logs: %v", err)
	}
	return len(records)
}

func TestIngestBatch_SingleTransaction(t *testing.T) {
	s, adapter := newTestStorage(t, 0)

	var batch []ingestRequest
	var replies []<-chan error
	for range 10 {
		req, reply := newTestRequest(t, "test", map[string]string{"level": "info"}, nil)
		batch = append(batch, req)
		replies = append(replies, reply)
	}

	s.ingestBatch(t.Context(), batch)

	for i, reply := range replies {
		if err := <-reply; err != nil {
			t.Fatalf("record %d: unexpected error: %v", i, err)
		}
	}
	if got := adapter.updates.Load(); got != 1 {
		t.Fatalf("expected 1 transaction, got %d", got)
	}
	if got := countRecords(t, s, "test"); got != 10 {
		t.Fatalf("expected 10 records, got %d", got)
	}
}

func TestIngestBatch_SplitsTooLargeBatch(t *testing.T) {
	// 8MB memtables allow about 1.2MB per transaction, 30 records of 60KB do not fit
	s, adapter := newTestStorage(t, 8<<20)

	var batch []ingestRequest
	var replies []<-chan error
	for range 30 {
		req, reply := newTestRequest(t, "test", map[string]string{"msg": strings.Repeat("x", 60_000)}, nil)
		batch = append(batch, req)
		replies = append(replies, reply)
	}

	s.ingestBatch(t.Context(), batch)

	for i, reply := range replies {
		if err := <-reply; err != nil {
			t.Fatalf("record %d: unexpected error: %v", i, err)
		}
	}
	if got := adapter.updates.Load(); got < 2 {
		t.Fatalf("expected the batch to be split, got %d transaction(s)", got)
	}
	if got := countRecords(t, s, "test"); got != 30 {
		t.Fatalf("expected 30 records, got %d", got)
	}
}

func TestIngestBatch_ReportsFailingRecordOnly(t *testing.T) {
	s, _ := newTestStorage(t, 0)

	var batch []ingestRequest
	var replies []<-chan error
	for i := range 5 {
		var val []byte
		if i == 2 {
			val = make([]byte, kv.MaxValueSize+1)
		}
		req, reply := newTestRequest(t, "test", map[string]string{"level": "info"}, val)
		batch = append(batch, req)
		replies = append(replies, reply)
	}

	s.ingestBatch(t.Context(), batch)

	for i, reply := range replies {
		err := <-reply
		if i == 2 {
			if !errors.Is(err, kv.ErrValueTooLarge) {
				t.Fatalf("record %d: expected ErrValueTooLarge, got %v", i, err)
			}
		} else if err != nil {
			t.Fatalf("record %d: unexpected error: %v", i, err)
		}
	}
	if got := countRecords(t, s, "test"); got != 4 {
		t.Fatalf("expected 4 records, got %d", got)
	}
}

func TestIngestBatch_SkipsCallersThatGaveUp(t *testing.T) {
	s, adapter := newTestStorage(t, 0)

	req, reply := newTestRequest(t, "test", map[string]string{"level": "info"}, nil)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	req.ctx = ctx

	s.ingestBatch(t.Context(), []ingestRequest{req})

	if err := <-reply; !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if got := adapter.updates.Load(); got != 0 {
		t.Fatalf("expected no transaction, got %d", got)
	}
}
