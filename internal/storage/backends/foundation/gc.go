package foundation

import (
	"context"
	"log/slog"
	"time"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/vladopajic/go-actor/actor"
)

// gcBatchSize bounds how many pairs are scanned per garbage-collection
// transaction so a single pass stays within FoundationDB's transaction limits.
const gcBatchSize = 1000

// gcWorker periodically deletes keys whose embedded TTL has expired.
//
// FoundationDB has no native TTL, so expired keys are otherwise only hidden
// lazily on read (see [FoundationQueryTx.Get] and the iterators); this worker
// reclaims their storage.
type gcWorker struct {
	adapter    *FoundationAdapter
	gcInterval time.Duration
}

var _ actor.Worker = (*gcWorker)(nil)

// NewGarbageCollector returns an [actor.Worker] that scans the adapter's
// subspace every gcInterval and removes expired keys. Errors are logged and do
// not stop the worker.
func NewGarbageCollector(adapter *FoundationAdapter, gcInterval time.Duration) actor.Worker {
	return &gcWorker{
		adapter:    adapter,
		gcInterval: gcInterval,
	}
}

func (w *gcWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case <-time.After(w.gcInterval):
		go func() {
			if err := collectExpired(ctx, w.adapter.db, w.adapter.sub); err != nil {
				slog.ErrorContext(
					ctx,
					"failed to collect expired keys",
					slog.String("channel", "storage.foundation"),
					slog.String("error", err.Error()),
				)
			}
		}()

		return actor.WorkerContinue
	}
}

// collectExpired scans the subspace and clears every key whose TTL has elapsed.
//
// Candidate keys are found with a snapshot read (which adds no read conflicts),
// then re-checked inside the deleting transaction with a plain read so a key a
// concurrent writer has just refreshed is never dropped.
func collectExpired(ctx context.Context, db fdb.Database, sub subspace.Subspace) error {
	beginKey, endKey := sub.FDBRangeKeys()
	begin := beginKey.FDBKey()
	end := endKey.FDBKey()

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		r := fdb.KeyRange{Begin: begin, End: end}

		result, err := db.ReadTransact(func(rtr fdb.ReadTransaction) (any, error) {
			if err := applyDeadline(ctx, rtr); err != nil {
				return nil, err
			}

			opts := fdb.RangeOptions{Limit: gcBatchSize, Mode: fdb.StreamingModeWantAll}
			return rtr.Snapshot().GetRange(r, opts).GetSliceWithError()
		})
		if err != nil {
			return err
		}

		pairs := result.([]fdb.KeyValue)
		if len(pairs) == 0 {
			return nil
		}

		var candidates []fdb.Key
		for _, pair := range pairs {
			if expired(decodeExpiresAt(pair.Value)) {
				candidates = append(candidates, pair.Key)
			}
		}

		if len(candidates) > 0 {
			_, err := db.Transact(func(tr fdb.Transaction) (any, error) {
				if err := applyDeadline(ctx, tr); err != nil {
					return nil, err
				}

				for _, key := range candidates {
					value, err := tr.Get(key).Get()
					if err != nil {
						return nil, err
					}

					if value != nil && expired(decodeExpiresAt(value)) {
						tr.Clear(key)
					}
				}

				return nil, nil
			})
			if err != nil {
				return err
			}
		}

		if len(pairs) < gcBatchSize {
			return nil
		}

		// Resume strictly after the last key scanned.
		lastKey := pairs[len(pairs)-1].Key
		next := make(fdb.Key, len(lastKey)+1)
		copy(next, lastKey)
		begin = next
	}
}
