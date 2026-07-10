# auth

The package at `internal/storage/databases/auth` implements the `AuthStorage`
contract from [interfaces](../../interfaces).

It owns the persistence of users, roles and personal access tokens, and the
credential and permission checks built on top of them, so the rest of FlowG can
reason about identity purely through the `AuthStorage` interface. The
implementation is backend-agnostic: it runs on top of any
[generic/kv](../../generic/kv) adapter.

## Responsibilities

- **Persistence** — stores users, roles and personal access tokens through a
  key-value adapter, keeping the auth key space isolated from the other domains.
- **Credentials** — verifies passwords and tokens, hashing secrets so plaintext
  is never persisted.
- **Authorization** — resolves the permission scopes granted to a user through
  its roles.
- **Snapshots** — satisfies `Streamable` so the auth database can be backed up
  and restored.

## Layout

- **storage.go** — the `Storage` type implementing `AuthStorage`, delegating each
  operation to the `transactions` subpackage inside a read or write transaction.
- **[transactions](transactions)** — the low-level read/write operations and the
  auth key-space layout.
