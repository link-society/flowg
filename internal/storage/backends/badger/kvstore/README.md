# kvstore

The package at `internal/storage/backends/badger/kvstore` is a
concurrency-safe wrapper around a [BadgerDB](https://github.com/dgraph-io/badger)
database.

It exists to give the BadgerDB-backed domain stores (`auth`, `config` and
`log`) a single, shared primitive for talking to the database. Rather than each
store opening and synchronizing its own database, they build on top of this
package, which serializes all access through an actor mailbox.

## Responsibilities

- **Serialized access** — `Storage` funnels every transaction through one actor
  worker, so concurrent callers never race on the underlying database and the
  domain stores stay free of locking code.
- **Transactions** — `View` runs a read-only transaction and `Update` runs a
  read-write one, each executing the caller-supplied function against a
  BadgerDB transaction.
- **Backup & restore** — `Backup` streams an incremental snapshot and `Restore`
  loads one back, underpinning FlowG's backup feature.
- **Wiring** — `NewStorage` returns an `fx` module that opens the database,
  starts the mailbox and binds their lifecycles to the application; `Options`
  and `DefaultOptions` configure it.

## Usage shape

A domain store provisions a `kvstore` through `NewStorage`, then issues
`View`/`Update` transactions to read and write its own key space:

```text
domain store ──▶ kvstore.View/Update ──▶ actor worker ──▶ BadgerDB
                 (issue transaction)      (serialize)      (persist)
```
