# log

The package at `internal/storage/databases/log` implements the `LogStorage`
contract from [interfaces](../../interfaces).

It persists ingested log records, organized into streams, along with each
stream's configuration, field set and inverted indices, and answers the
time-ranged, filtered queries the query engine issues. It exposes everything
through the `LogStorage` interface, so consumers never depend on a concrete
database. The implementation is backend-agnostic: it runs on top of any
[generic/kv](../../generic/kv) adapter.

## Responsibilities

- **Ingestion** — stores log records under time-ordered keys with a per-stream
  retention TTL, tracking each stream's field set as records arrive.
- **Indexing** — maintains inverted indices for the fields a stream is configured
  to index, back-filling and dropping them as the configuration changes, so
  queries can narrow candidates without scanning every record.
- **Querying** — returns the records of a stream within a time range that satisfy
  a filter and the requested indexed-field constraints.
- **Retention** — enforces each stream's size budget through a background garbage
  collector in addition to the per-record TTLs.
- **Snapshots** — satisfies `Streamable` so the log database can be backed up and
  restored.

## Layout

- **storage.go** — the `Storage` type implementing `LogStorage`, delegating each
  operation to the `transactions` subpackage inside a read or write transaction.
- **gc.go** — `NewGarbageCollector`, the background worker that periodically
  enforces stream retention-size budgets.
- **[transactions](transactions)** — the low-level read/write operations and the
  log key-space layout.
