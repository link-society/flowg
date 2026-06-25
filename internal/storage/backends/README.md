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
