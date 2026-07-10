# backends

The packages under `internal/storage/backends` provide the concrete backends
that satisfy the storage contracts declared in [interfaces](../interfaces).

A backend supplies a [kv.Adapter](../generic/kv) for a particular database
technology and wires it to the backend-agnostic domain stores under
[databases](../databases), exposing them as the `AuthStorage`, `ConfigStorage`
and `LogStorage` interfaces. Keeping them here, behind the interfaces, lets FlowG
adopt a different storage technology without changing any of its consumers.

## Layout

- **[badger](badger)** — the default backend, built on
  [BadgerDB](https://github.com/dgraph-io/badger). It provides a BadgerDB-backed
  key-value adapter and per-domain `fx` modules that combine it with the domain
  stores.
- **[foundation](foundation)** — third-party backend, built on
  [FoundationDB](https://www.foundationdb.org/).
