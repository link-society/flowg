# backends

The packages under `internal/storage/backends` provide the concrete
implementations of the storage contracts declared in
[internal/storage](../).

Each backend translates the domain interfaces — `AuthStorage`, `ConfigStorage`
and `LogStorage` — into operations against a particular database technology.
Keeping them here, behind the interfaces, lets FlowG adopt a different storage
technology without changing any of its consumers.

## Layout

- **[badger](badger)** — the default backend, built on
  [BadgerDB](https://github.com/dgraph-io/badger). It contains one package per
  domain interface plus the shared `kvstore` that wraps the embedded key-value
  database they all build upon.
- **[foundationdb](foundationdb)** — an optional backend built on
  [FoundationDB](https://www.foundationdb.org), a distributed key-value store.
  It replaces BadgerDB for clustered deployments and follows the same package
  structure — one sub-package per domain interface under
  `foundationdb/concrete/`. Each sub-package exposes `DefaultOptions()` and
  `NewStorage()` following the same conventions as the badger backend.
