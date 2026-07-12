# kv

The package at `internal/storage/generic/kv` defines the generic key-value store
abstraction that FlowG's domain storage is built on.

It models a database as a set of composite keys mapped to byte values, accessed
through transactions, without committing to any particular engine. The domain
implementations under [databases](../../databases) depend only on this package,
so the same logic runs on top of any backend that provides an `Adapter`.

## Concepts

- **`Key` / `KeyRange` / `KeySlice`** — a key is an ordered list of string
  segments (`Key`), joined by the backend into its native key format. `KeyRange`
  bounds an iteration (`From` and `To` are both inclusive of their subtree) and
  `KeySlice` is a sortable sequence of keys.
- **`Value` / `Pair`** — `Value` is an arbitrary byte payload; `Pair` exposes a
  stored key together with its value, size estimate and expiration time during
  iteration.
- **`QueryTx` / `MutationTx`** — the transaction contracts. `QueryTx` offers
  read-only access (`Get`, `IterKeys`, `IterPairs`); `MutationTx` extends it with
  writes (`Set`, `SetWithTTL`, `Clear`).
- **`Adapter`** — the entry point to a backend. It runs read-only (`View`) and
  read-write (`Update`) transactions and streams incremental snapshots with
  `Backup` and `Restore`. It is parametrized by the concrete transaction types so
  backends can expose their own transaction implementations.

## Limits

Backends adopt FoundationDB's hard size limits as part of the contract, so a
mutation that would be rejected by one engine is rejected by all of them up
front rather than failing deep inside a commit:

- **`MaxKeySize`** (10 KB) — the maximum size of an encoded key. Every key-
  addressed operation (`Get`, `Set`, `SetWithTTL`, `Clear`) returns
  `ErrKeyTooLarge` when exceeded, so an oversized key can neither be written,
  read nor deleted.
- **`MaxValueSize`** (100 KB) — the maximum size of an encoded value (including
  any envelope a backend adds, such as FoundationDB's expiry header). `Set` /
  `SetWithTTL` return `ErrValueTooLarge` when exceeded.

Backends validate the size of the key and value **as they will be stored** using
`CheckKeySize` / `CheckValueSize`.

Consumers may treat `ErrKeyTooLarge` / `ErrValueTooLarge` as a signal to degrade
gracefully rather than fail: the log storage, for instance, skips indexing a
field value that would overflow an index key instead of rejecting the record
(see [databases/log/transactions](../../databases/log/transactions)).


## Layout

- **types.go** — `Key`, `KeyRange`, `KeySlice`, `Value` and the `Pair` contract.
- **txn.go** — the `QueryTx` and `MutationTx` transaction contracts.
- **adapter.go** — the `Adapter` contract binding transactions to a backend and
  exposing backup/restore.
- **errors.go** — shared sentinel errors such as `ErrNotSupported`,
  `ErrKeyTooLarge` and `ErrValueTooLarge`.
- **limits.go** — the `MaxKeySize` / `MaxValueSize` limits and the
  `CheckKeySize` / `CheckValueSize` helpers backends use to enforce them.