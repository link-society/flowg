package logstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/vladopajic/go-actor/actor"
)

type indexerMessage struct {
	Stream  string
	Field   string
	Indexed bool
}

type indexer struct {
	mbox  actor.MailboxSender[indexerMessage]
	actor actor.Actor
}

type indexerWorker struct {
	mbox actor.MailboxReceiver[indexerMessage]
	db   *badger.DB
}

func newIndexer(db *badger.DB) *indexer {
	mbox := actor.NewMailbox[indexerMessage]()
	worker := actor.New(&indexerWorker{mbox: mbox, db: db})

	return &indexer{
		mbox:  mbox,
		actor: actor.Combine(mbox, worker).Build(),
	}
}

func (idx *indexer) Start() {
	idx.actor.Start()
}

func (idx *indexer) Stop() {
	idx.actor.Stop()
}

func (idx *indexer) IndexField(stream, field string) {
	idx.mbox.Send(
		context.Background(),
		indexerMessage{Stream: stream, Field: field, Indexed: true},
	)
}

func (idx *indexer) UnindexField(stream, field string) {
	idx.mbox.Send(
		context.Background(),
		indexerMessage{Stream: stream, Field: field, Indexed: false},
	)
}

func (w *indexerWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case msg, ok := <-w.mbox.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		if msg.Indexed {
			w.indexField(ctx, msg.Stream, msg.Field)
		} else {
			w.unindexField(ctx, msg.Stream, msg.Field)
		}

		return actor.WorkerContinue
	}
}

func (w *indexerWorker) indexField(ctx actor.Context, stream, field string) {
	slog.DebugContext(
		ctx,
		"Indexing field",
		"channel", "logstorage",
		"stream", stream,
		"field", field,
	)

	err := w.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.Prefix = []byte("entry:" + stream + ":")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.KeyCopy(nil)

			var entry LogEntry
			err := item.Value(func(val []byte) error {
				if err := json.Unmarshal(val, &entry); err != nil {
					return fmt.Errorf("could not unmarshal log entry '%s': %w", key, err)
				}

				return nil
			})
			if err != nil {
				return err
			}

			ts := item.ExpiresAt()
			retentionTime := int64(0)
			if ts != 0 {
				retentionTime = int64(ts) - time.Now().Unix()
			}

			index := newFieldIndex(txn, stream, field, entry.Fields[field])
			if err := index.AddKey(key, retentionTime); err != nil {
				return fmt.Errorf(
					"could not index field '%s' for entry '%s': %w",
					field, key, err,
				)
			}
		}

		return nil
	})

	if err != nil {
		slog.ErrorContext(
			ctx,
			"Could not index field",
			"channel", "logstorage",
			"stream", stream,
			"field", field,
			"error", err.Error(),
		)
	}
}

func (w *indexerWorker) unindexField(ctx actor.Context, stream, field string) {
	slog.DebugContext(
		ctx,
		"Unindexing field",
		"channel", "logstorage",
		"stream", stream,
		"field", field,
	)

	err := w.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(fmt.Sprintf("index:%s:field:%s:", stream, field))
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			indexKey := it.Item().KeyCopy(nil)

			if err := txn.Delete(indexKey); err != nil {
				return fmt.Errorf("could not delete index key '%s': %w", indexKey, err)
			}
		}

		return nil
	})

	if err != nil {
		slog.ErrorContext(
			ctx,
			"Could not unindex field",
			"channel", "logstorage",
			"stream", stream,
			"field", field,
			"error", err.Error(),
		)
	}
}
