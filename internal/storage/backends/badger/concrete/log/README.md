# log

The package at `internal/storage/backends/badger/log` implements the
`LogStorage` contract from [internal/storage](../../../../) on top of
[BadgerDB](https://github.com/dgraph-io/badger).

It is the default log backend. It owns the persistence, indexing and querying
of ingested log records, exposing them through the `LogStorage` interface so the
ingestion pipeline and query API never depend on BadgerDB directly.

## Responsibilities

- **Ingestion** — appends log records to a stream, persisting them in a
  dedicated [kvstore](../../kvstore).
- **Indexing** — maintains per-field indices so records can be looked up by
  their fields rather than scanned in full.
- **Querying** — answers bounded, filtered searches over a stream's records.
- **Retention** — reclaims space for expired records through its garbage
  collection routine.
- **Snapshots** — satisfies `Streamable` so the log database can be backed up
  and restored.
- **Wiring** — `NewStorage` returns an `fx` module providing a `LogStorage`;
  `Options` and `DefaultOptions` configure where and how the database is opened.

## Layout

- **main.go** — the `LogStorage` implementation and its `fx` wiring.
- **gc.go** — the background garbage collection of expired records.
- **[transactions/](transactions)** — the low-level read/write and indexing operations against
  the key space.
