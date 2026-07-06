# foundation

The packages under `internal/storage/backends/foundation` implement FlowG's
storage contracts on top of [FoundationDB](https://foundationdb.org).

Each domain interface declared in [internal/storage](../../) has its own package
here, and they all share a single key-value primitive so they behave consistently
and can be managed through the same lifecycle.

## Layout

- **[kvstore](kvstore)** — the shared, concurrency-safe wrapper around a
  FoundationDB database. Every domain store is built on top of it.
- **[concrete](concrete)** — the per-domain implementations of the storage
  interfaces:
  - **[auth](concrete/auth)** — implements `AuthStorage`: users, roles, personal
    access tokens and the password/permission checks around them.
  - **[config](concrete/config)** — implements `ConfigStorage`: pipelines,
    transformers, forwarders and system configuration.
  - **[log](concrete/log)** — implements `LogStorage`: ingestion, indexing and
    querying of log records.

## Design

Each domain package exposes the same shape: an `Options` struct, a
`DefaultOptions` constructor and a `NewStorage` function that returns an `fx`
module. `NewStorage` provisions a dedicated `kvstore` and adapts its
transactions into the domain interface, so the rest of the application only ever
sees the interface from [internal/storage](../../), never FoundationDB itself.
