# config

The package at `internal/storage/databases/config` implements the
`ConfigStorage` contract from [interfaces](../../interfaces).

It persists the resources that define how FlowG processes logs — pipelines,
transformers and forwarders — together with the system-wide configuration,
exposing them through the `ConfigStorage` interface so the engines and API never
depend on a concrete database. The implementation is backend-agnostic: it runs on
top of any [generic/kv](../../generic/kv) adapter.

## Responsibilities

- **Resource persistence** — stores and retrieves pipelines, transformers and
  forwarders through a key-value adapter, migrating them to the latest model
  version on read when needed.
- **System configuration** — persists and caches global settings such as the
  allowed origins used by the ingestion endpoints, validating them on write.
- **Snapshots** — satisfies `Streamable` so the configuration database can be
  backed up and restored.

## Layout

- **storage.go** — the `Storage` type implementing `ConfigStorage`, delegating
  each operation to the `transactions` subpackage inside a read or write
  transaction.
- **[transactions](transactions)** — the low-level read/write operations and the
  config key-space layout.
