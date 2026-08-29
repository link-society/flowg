package foundation

import (
	"context"
	"fmt"
	"io"
	"time"

	"crypto/rand"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/storage/generic/kv"
)

// FoundationApiVersion is the FoundationDB API version negotiated with the
// installed client library.
const FoundationApiVersion = 730

// defaultGCInterval is how often the expired-key collector runs when
// AdapterOptions.GCInterval is left unset.
const defaultGCInterval = 5 * time.Minute

// ChangeLogKey is the key bumped by every Update transaction; watching it (see
// [FoundationAdapter.Watch]) observes every change made through the adapter,
// from any node of the cluster.
var ChangeLogKey = kv.Key{"meta", "changelog"}

// AdapterOptions configures how a FoundationDB-backed [kv.Adapter] connects to
// its cluster and where within the key space it stores its data.
type AdapterOptions struct {
	// LogChannel names the logging channel and the fx module the store lives in.
	LogChannel string
	// ClusterFile is the path to the fdb.cluster file; empty uses the default.
	ClusterFile string
	// ConnectionString is the FoundationDB connection string; empty uses the default.
	ConnectionString string
	// KeySpace is the root prefix shared by every FlowG storage (e.g. "flowg").
	KeySpace string
	// Namespace is the subspace this storage owns (e.g. "config", "auth", "log").
	Namespace string
	// GCInterval is how often expired keys are swept; zero uses defaultGCInterval.
	GCInterval time.Duration
}

// FoundationAdapter is a [kv.Adapter] backed by FoundationDB.
//
// Every operation is scoped to the subspace KeySpace/Namespace. That prefix is
// applied on write and stripped on read, so consumers only ever observe logical
// [kv.Key]s.
type FoundationAdapter struct {
	db  fdb.Database
	sub subspace.Subspace
}

var _ kv.Adapter[*FoundationQueryTx, *FoundationMutationTx] = (*FoundationAdapter)(nil)

// NewAdapter builds an fx module that provides a [kv.Adapter] backed by
// FoundationDB.
//
// The resulting [kv.Adapter] is published under the name given by
// AdapterOptions.LogChannel so several adapters can coexist in the same
// container.
func NewAdapter(opts AdapterOptions) fx.Option {
	makeAdapter := func(lc fx.Lifecycle) (*FoundationAdapter, error) {
		if err := fdb.APIVersion(FoundationApiVersion); err != nil {
			return nil, fmt.Errorf("failed to select FoundationDB API version: %w", err)
		}

		var (
			db  fdb.Database
			err error
		)
		if opts.ConnectionString != "" {
			db, err = fdb.OpenWithConnectionString(opts.ConnectionString)
		} else if opts.ClusterFile != "" {
			db, err = fdb.OpenDatabase(opts.ClusterFile)
		} else {
			err = fmt.Errorf("either ConnectionString or ClusterFile must be provided")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}

		adapter := &FoundationAdapter{
			db:  db,
			sub: subspace.Sub(opts.KeySpace).Sub(opts.Namespace),
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				db.Close()
				return nil
			},
		})

		gcInterval := opts.GCInterval
		if gcInterval <= 0 {
			gcInterval = defaultGCInterval
		}

		gc := actor.New(NewGarbageCollector(adapter, gcInterval))
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				gc.Start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				gc.Stop()
				return nil
			},
		})

		return adapter, nil
	}

	module := fmt.Sprintf("kv.adapter.%s", opts.LogChannel)

	return fx.Module(
		module,
		fx.Provide(
			fx.Annotate(
				makeAdapter,
				fx.ResultTags(fmt.Sprintf(`name:"%s"`, opts.LogChannel)),
			),
		),
	)
}

// View implements [kv.Adapter.View]. It runs inside a read-only FoundationDB
// transaction, which is never committed; FoundationDB retries automatically on
// retryable errors.
func (a *FoundationAdapter) View(ctx context.Context, txnFn func(txn *FoundationQueryTx) error) error {
	_, err := a.db.ReadTransact(func(rtr fdb.ReadTransaction) (any, error) {
		if err := applyDeadline(ctx, rtr); err != nil {
			return nil, err
		}

		return nil, txnFn(&FoundationQueryTx{concrete: rtr, sub: a.sub})
	})
	return err
}

// Update implements [kv.Adapter.Update]. It runs inside a read-write
// FoundationDB transaction, which is committed on success; FoundationDB retries
// automatically on conflict. [ChangeLogKey] is set to a fresh random value in
// the same transaction, so watchers observe every mutation.
func (a *FoundationAdapter) Update(ctx context.Context, txnFn func(txn *FoundationMutationTx) error) error {
	_, err := a.db.Transact(func(tr fdb.Transaction) (any, error) {
		if err := applyDeadline(ctx, tr); err != nil {
			return nil, err
		}

		if err := txnFn(&FoundationMutationTx{concrete: tr, sub: a.sub}); err != nil {
			return nil, err
		}

		value := make([]byte, 16)
		rand.Read(value)
		tr.Set(a.sub.Pack(keyToTuple(ChangeLogKey)), value)

		return nil, nil
	})
	return err
}

// Watch arms a FoundationDB key watch and returns a channel that fires once
// when the key's value changes. Cancelling ctx cancels the watch, which then
// fires with the cancellation error.
func (a *FoundationAdapter) Watch(ctx context.Context, key kv.Key) (<-chan error, error) {
	var fut fdb.FutureNil
	_, err := a.db.Transact(func(tr fdb.Transaction) (any, error) {
		if err := applyDeadline(ctx, tr); err != nil {
			return nil, err
		}

		fut = tr.Watch(a.sub.Pack(keyToTuple(key)))
		return nil, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to arm watch: %w", err)
	}

	firedC := make(chan error, 1)
	resolvedC := make(chan struct{})

	go func() {
		firedC <- fut.Get()
		close(resolvedC)
	}()

	go func() {
		select {
		case <-ctx.Done():
			fut.Cancel()
		case <-resolvedC:
		}
	}()

	return firedC, nil
}

// Backup implements [kv.Adapter.Backup].
//
// FoundationDB exposes no snapshot streaming primitive through the client API;
// backups are taken out-of-band with the fdbbackup tooling.
func (a *FoundationAdapter) Backup(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	return 0, kv.ErrNotSupported
}

// Restore implements [kv.Adapter.Restore].
//
// FoundationDB exposes no bulk load primitive through the client API; restores
// are performed out-of-band with the fdbrestore tooling.
func (a *FoundationAdapter) Restore(ctx context.Context, r io.Reader) error {
	return kv.ErrNotSupported
}

// deadlineOptioner is satisfied by both [fdb.Transaction] and
// [fdb.ReadTransaction], exposing the transaction options.
type deadlineOptioner interface {
	Options() fdb.TransactionOptions
}

// applyDeadline aborts early when ctx is already done and, when ctx carries a
// deadline, maps it onto the transaction's timeout so it cannot outlive the
// caller's expectations.
func applyDeadline(ctx context.Context, tr deadlineOptioner) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline).Milliseconds()
		if remaining <= 0 {
			return context.DeadlineExceeded
		}

		if err := tr.Options().SetTimeout(remaining); err != nil {
			return err
		}
	}

	return nil
}
