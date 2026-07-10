# interfaces

The package at `internal/storage/interfaces` declares the storage contracts that
the rest of FlowG depends on.

It exists to decouple the application from any particular database technology.
Engines, services and API operations are written against the interfaces declared
here; the concrete backends that satisfy them live under
[backends](../backends), and the reusable implementations they build upon live
under [databases](../databases). Swapping or adding a backend therefore never
requires touching the consumers.

## Contracts

- **`AuthStorage`** — persistence of identities and permissions: roles, users,
  their permission scopes, and the personal access tokens used to authenticate
  API calls.
- **`ConfigStorage`** — persistence of the resources that define how FlowG
  processes logs: transformers, pipelines, forwarders, and the global system
  configuration.
- **`LogStorage`** — persistence and querying of ingested log records, organized
  into streams, together with each stream's configuration and field indices.
- **`Streamable`** — the snapshot/restore capability embedded by every storage
  interface. It exposes `Dump` and `Load` so any backend can be backed up and
  restored, powering the backup and restore features.

## Layout

- **auth.go** — the `AuthStorage` contract.
- **config.go** — the `ConfigStorage` contract.
- **log.go** — the `LogStorage` contract.
- **streamable.go** — the `Streamable` snapshot/restore contract embedded by the
  three domain interfaces.
