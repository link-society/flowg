# badger

The packages under `internal/storage/backends/badger` implement FlowG's storage
contracts on top of [BadgerDB](https://github.com/dgraph-io/badger), an
embedded key-value database.

They are the default backend wired into the application. Each domain interface
declared in [internal/storage](../../) has its own package here, and they all
share a single key-value primitive so they behave consistently and can be
managed through the same lifecycle.

## Layout

- **[kvstore](kvstore)** — the shared, concurrency-safe wrapper around a
  BadgerDB database. Every domain store is built on top of it.
- **[auth](auth)** — implements `AuthStorage`: users, roles, personal access
  tokens and the password/permission checks around them.
- **[config](config)** — implements `ConfigStorage`: pipelines, transformers,
  forwarders and system configuration.
- **[log](log)** — implements `LogStorage`: ingestion, indexing and querying of
  log records.

## Design

Each domain package exposes the same shape: an `Options` struct, a
`DefaultOptions` constructor and a `NewStorage` function that returns an `fx`
module. `NewStorage` provisions a dedicated `kvstore` and adapts its
transactions into the domain interface, so the rest of the application only ever
sees the interface from [internal/storage](../../), never BadgerDB itself.
