# config

The package at `internal/storage/backends/badger/config` implements the
`ConfigStorage` contract from [internal/storage](../../../../) on top of
[BadgerDB](https://github.com/dgraph-io/badger).

It is the default configuration backend. It persists the resources that define
how FlowG processes logs — pipelines, transformers and forwarders — together
with the system-wide configuration, exposing them through the `ConfigStorage`
interface so the engines and API never depend on BadgerDB directly.

## Responsibilities

- **Resource persistence** — stores and retrieves pipelines, transformers and
  forwarders in a dedicated [kvstore](../../kvstore).
- **System configuration** — persists global settings such as the allowed
  origins used by the ingestion endpoints.
- **Snapshots** — satisfies `Streamable` so the configuration database can be
  backed up and restored.
- **Wiring** — `NewStorage` returns an `fx` module providing a `ConfigStorage`;
  `Options` and `DefaultOptions` configure where and how the database is opened.

## Layout

- **main.go** — the `ConfigStorage` implementation and its `fx` wiring.
- **migrator.go** — schema migrations applied when the database is opened.
- **[transactions/](transactions)** — the low-level read/write operations against the key space.
