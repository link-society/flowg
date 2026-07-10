# databases

The packages under `internal/storage/databases` provide backend-agnostic
implementations of the storage contracts declared in
[interfaces](../interfaces).

Each domain store is written once, against the generic key-value abstraction in
[generic/kv](../generic/kv), and runs on top of any backend that provides a
`kv.Adapter`. The concrete backends under [backends](../backends) simply
instantiate these stores with their own adapter, so the domain logic — key
layout, indexing, migrations, reconciliation — never has to be reimplemented per
backend.

## Layout

- **[auth](auth)** — implements `AuthStorage`: users, roles, personal access
  tokens and the password/permission checks around them.
- **[config](config)** — implements `ConfigStorage`: pipelines, transformers,
  forwarders and system configuration.
- **[log](log)** — implements `LogStorage`: ingestion, indexing, querying and
  retention of log records.

## Design

Every domain package follows the same shape: a generic `Storage[QTx, MTx]` type
holding a `kv.Adapter`, a `NewStorage` constructor, and a set of methods that
open a transaction on the adapter and delegate the actual key manipulation to a
`transactions` subpackage. The `transactions` subpackages own the key space
layout and the low-level read/write logic; see their own READMEs for the details.
