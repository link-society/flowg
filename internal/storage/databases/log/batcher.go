package log

import (
	"context"

	"github.com/vladopajic/go-actor/actor"

	"go.link-society.com/flowg/pkg/models"

	"go.link-society.com/flowg/internal/storage/generic/kv"

	"go.link-society.com/flowg/internal/storage/databases/log/transactions"
)

// maxIngestBatch bounds how many records a single ingestion transaction stores.
const maxIngestBatch = 256

// ingestRequest is an [Storage.Ingest] call waiting for the ingest batcher.
type ingestRequest struct {
	ctx     context.Context
	stream  string
	record  *models.LogRecord
	key     kv.Key
	val     []byte
	replyTo chan<- error
}

type ingestBatcher[QTx kv.QueryTx, MTx kv.MutationTx] struct {
	storage *Storage[QTx, MTx]
}

var _ actor.Worker = (*ingestBatcher[kv.QueryTx, kv.MutationTx])(nil)

// NewIngestBatcher returns an [actor.Worker] that stores the records of pending
// [Storage.Ingest] calls, up to maxIngestBatch at a time, in a single
// transaction. It must run for Ingest calls to complete.
func NewIngestBatcher[QTx kv.QueryTx, MTx kv.MutationTx](storage *Storage[QTx, MTx]) actor.Worker {
	return &ingestBatcher[QTx, MTx]{storage: storage}
}

func (w *ingestBatcher[QTx, MTx]) DoWork(ctx actor.Context) actor.WorkerStatus {
	var batch []ingestRequest

	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case req := <-w.storage.ingestC:
		batch = append(batch, req)
	}

	// take the calls already waiting, without waiting for more
drain:
	for len(batch) < maxIngestBatch {
		select {
		case req := <-w.storage.ingestC:
			batch = append(batch, req)
		default:
			break drain
		}
	}

	w.storage.ingestBatch(ctx, batch)

	return actor.WorkerContinue
}

// ingestBatch stores the records of batch in one transaction, and replies to
// every caller with its outcome. Callers that gave up are skipped.
func (s *Storage[QTx, MTx]) ingestBatch(ctx context.Context, batch []ingestRequest) {
	pending := batch[:0]
	for _, req := range batch {
		if err := req.ctx.Err(); err != nil {
			req.replyTo <- err
			continue
		}
		pending = append(pending, req)
	}

	if len(pending) == 0 {
		return
	}

	err := s.adapter.Update(ctx, func(txn MTx) error {
		for _, req := range pending {
			err := transactions.Ingest(s.logger, txn, req.stream, req.record, req.key, req.val, s.streamConfigOverride)
			if err != nil {
				return err
			}
		}
		return nil
	})

	for _, req := range pending {
		req.replyTo <- err
	}
}
