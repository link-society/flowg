# foundation

The packages under `internal/storage/backends/foundation` implement FlowG's
storage contracts on top of [FoundationDB](https://www.foundationdb.org/), a
distributed, transactional key-value database.

It provides a FoundationDB-backed [kv.Adapter](../../generic/kv), which the
backend-agnostic domain stores under [databases](../../databases) run on
unchanged. The per-domain packages under [concrete](concrete) assemble the
adapter and the matching domain store into ready-to-use `fx` modules.

> **Note:** the Go bindings are cgo wrappers around the FoundationDB client
> library, so building this backend requires `libfdb_c` and its headers (the
> `foundationdb-clients` package) to be installed.

## Design

`NewAdapter` returns an `fx` module that provides a `FoundationAdapter`, a
`kv.Adapter` implementation. `View` runs inside a read-only transaction
(`db.ReadTransact`, never committed) and `Update` inside a read-write
transaction (`db.Transact`, committed on success); FoundationDB retries both
automatically on retryable errors and conflicts.

### Key space

The whole backend lives under a subspace named `KeySpace/Namespace` (e.g.
`flowg/config`), where `KeySpace` is shared by every FlowG storage and
`Namespace` names this particular storage. A composite `kv.Key` is encoded as a
FoundationDB tuple and packed into that subspace; on the way out the subspace
prefix is stripped again, so consumers of the `kv.Adapter` only ever observe
their own logical keys and never see the prefix.

### Time-to-live

FoundationDB has no native TTL. `SetWithTTL` stores the expiration timestamp in
an 8-byte header prepended to the value, and reads (`Get` and the iterators)
skip entries whose header has elapsed. Because expired keys are only hidden
lazily on read, a background garbage collector — started automatically with the
adapter's `fx` lifecycle — periodically scans the subspace and physically clears
them.

### Backup & restore

`Backup` and `Restore` return `kv.ErrNotSupported`: FoundationDB exposes no
snapshot streaming or bulk load primitive through the client API. Backups are
taken out-of-band with the `fdbbackup` / `fdbrestore` tooling instead.

## Layout

- **adapter.go** — `FoundationAdapter` and `NewAdapter`, plus `AdapterOptions`.
  The `kv.Adapter` and the `fx` module that opens the database, wires its
  lifecycle and starts the garbage collector.
- **txn.go** — `FoundationQueryTx` and `FoundationMutationTx`, the read-only and
  read-write transaction adapters implementing `kv.QueryTx` and `kv.MutationTx`.
- **types.go** — the key/tuple encoding, the value envelope carrying the TTL, and
  the `kv.Pair` adapter over FoundationDB key-value pairs.
- **gc.go** — `NewGarbageCollector`, the background worker that periodically
  deletes keys whose embedded TTL has expired.
- **[concrete](concrete)** — the per-domain `fx` modules that combine the adapter
  with a domain store:
  - **[auth](concrete/auth)** — provides `AuthStorage`.
  - **[config](concrete/config)** — provides `ConfigStorage`.
  - **[log](concrete/log)** — provides `LogStorage` and its retention garbage
    collector.
