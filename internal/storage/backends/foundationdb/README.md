# foundationdb

The FoundationDB backend implements the `AuthStorage`, `ConfigStorage` and
`LogStorage` interfaces declared in `internal/storage` on top of
[FoundationDB](https://www.foundationdb.org), a distributed key-value store
with built-in transactions, sharding and replication.

## Layout

- **[kvstore](kvstore)** — a concurrency-safe FoundationDB wrapper that
  provides `View` (read-only) and `Update` (read-write) transaction methods
  along with custom `Backup` and `Restore` operations. All access goes through
  an actor mailbox, matching the same pattern used by the BadgerDB backend.
- **[concrete/auth](concrete/auth)** — `AuthStorage` implementation using
  FDB subspaces. Handles roles, users, personal access tokens, password
  verification and permission checks. Password hashing uses Argon2id.
- **[concrete/config](concrete/config)** — `ConfigStorage` implementation
  for pipelines, transformers, forwarders and system configuration.
- **[concrete/log](concrete/log)** — `LogStorage` implementation with log
  ingestion, field indexing, time-range queries, and a background garbage
  collection worker (since FoundationDB does not support key-level TTL).

## Usage

Select the FoundationDB backend through the configuration file or the
`FLOWG_STORAGE_BACKEND` environment variable:

```hcl
storage {
  backend "foundationdb" {
    connection_string = ""
  }
}
```

An empty `connection_string` tells the client to use the default FDB cluster
file (`/etc/foundationdb/fdb.cluster`). When set, it overrides this path.

## Key Space

The backend follows the same logical key space as the BadgerDB backend,
documented in `website/docs/design/storage.md` and
`website/docs/design/auth.md`, but encodes keys using FDB subspaces and tuples
instead of flat colon-delimited strings.

## Known Limitations

- **Backup/Restore**: FDB's native backup/restore is not exposed through the
  Go client library. The `Dump`/`Load` methods fall back to iterating over
  all keys in the namespace and serializing them as newline-delimited JSON.
- **TTL**: FoundationDB does not support key-level time-to-live. Log
  retention relies entirely on the garbage collection worker.
- **Dependencies**: The backend requires the FDB C client library
  (`libfdb_c`) and CGO. It cannot be compiled with `CGO_ENABLED=0`.
