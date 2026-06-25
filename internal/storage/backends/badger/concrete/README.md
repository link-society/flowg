# concrete

The packages under `internal/storage/backends/badger/concrete` are the
per-domain implementations of FlowG's storage contracts on top of
[BadgerDB](https://github.com/dgraph-io/badger).

They are grouped here to separate the concrete domain stores from the shared
BadgerDB primitives that sit alongside them (the
[kvstore](../kvstore) wrapper and the `BadgerLogger`). Each package adapts the
shared key-value store into one of the interfaces declared in
[internal/storage](../../../).

## Layout

- **[auth](auth)** — implements `AuthStorage`: users, roles, personal access
  tokens and the password/permission checks around them.
- **[config](config)** — implements `ConfigStorage`: pipelines, transformers,
  forwarders and system configuration.
- **[log](log)** — implements `LogStorage`: ingestion, indexing and querying of
  log records.
