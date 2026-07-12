# badger

The packages under `internal/storage/backends/badger` implement FlowG's storage
contracts on top of [BadgerDB](https://github.com/dgraph-io/badger), an
embedded key-value database.

This is the default backend wired into the application. It provides a single
BadgerDB-backed [kv.Adapter](../../generic/kv), which the backend-agnostic domain
stores under [databases](../../databases) run on unchanged. The per-domain
packages under [concrete](concrete) assemble the adapter and the matching domain
store into ready-to-use `fx` modules.

## Design

`NewAdapter` returns an `fx` module that provides a `BadgerAdapter`, a
`kv.Adapter` implementation. All database access is serialized through an actor
mailbox: `View`, `Update`, `Backup` and `Restore` each enqueue an operation that
a single worker goroutine applies to the BadgerDB handle, and write transactions
are transparently retried on conflict. Composite `kv.Key`s are joined into
BadgerDB's flat key space with a reserved separator byte — the ASCII `ESC`
control byte (`0x1B`, the `keySeparator` constant in `types.go`) — which the
segment values FlowG stores (stream names, field names, item names) never
contain, so the join/split is lossless. A `BadgerTx` adapts a BadgerDB
transaction to the `kv.QueryTx` / `kv.MutationTx` contracts.

## Layout

- **adapter.go** — `BadgerAdapter` and `NewAdapter`, plus `AdapterOptions`. The
  actor-based `kv.Adapter` and the `fx` module that wires its lifecycle.
- **txn.go** — `BadgerTx`, the transaction adapter implementing `kv.QueryTx` and
  `kv.MutationTx`.
- **types.go** — the key encoding and the `kv.Pair` adapter over BadgerDB items.
- **messages.go** — the internal operations (view, update, backup, restore) sent
  through the actor mailbox.
- **logging.go** — the adapter that routes BadgerDB's internal logs into FlowG's
  `slog` output.
- **[concrete](concrete)** — the per-domain `fx` modules that combine the adapter
  with a domain store (and its startup migrations):
  - **[auth](concrete/auth)** — provides `AuthStorage`.
  - **[config](concrete/config)** — provides `ConfigStorage`.
  - **[log](concrete/log)** — provides `LogStorage` and its retention garbage
    collector.
