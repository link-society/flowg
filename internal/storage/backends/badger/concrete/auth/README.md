# auth

The package at `internal/storage/backends/badger/auth` implements the
`AuthStorage` contract from [internal/storage](../../../../) on top of
[BadgerDB](https://github.com/dgraph-io/badger).

It is the default authentication and authorization backend. It owns the
persistence of users, roles and personal access tokens, and the credential and
permission checks built on top of them, so the rest of FlowG can reason about
identity purely through the `AuthStorage` interface.

## Responsibilities

- **Persistence** — stores users, roles and personal access tokens in a
  dedicated [kvstore](../../kvstore), keeping the auth key space isolated from the
  other domains.
- **Credentials** — verifies passwords and tokens, hashing secrets so plaintext
  is never persisted.
- **Authorization** — resolves the permission scopes granted to a user through
  its roles.
- **Snapshots** — satisfies `Streamable` so the auth database can be backed up
  and restored.
- **Wiring** — `NewStorage` returns an `fx` module providing an `AuthStorage`;
  `Options` and `DefaultOptions` configure where and how the database is opened.

## Layout

- **main.go** — the `AuthStorage` implementation and its `fx` wiring.
- **migrator.go** — schema migrations applied when the database is opened.
- **[transactions/](transactions)** — the low-level read/write operations against the key space.
