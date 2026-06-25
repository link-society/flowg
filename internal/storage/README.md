# storage

The package at `internal/storage` defines the storage contracts that the rest
of FlowG depends on.

It exists to decouple the application from any particular database technology.
The engines, services and API operations are written against the interfaces
declared here; the concrete backends that satisfy them live under
[backends](backends) and are wired in at startup. Swapping or adding a backend
therefore never requires touching the consumers.

## Responsibilities

- **Domain contracts** — `AuthStorage`, `ConfigStorage` and `LogStorage`
  describe the operations FlowG needs to manage, respectively, identities and
  permissions, pipeline/transformer/forwarder configuration, and ingested log
  records.
- **Shared capability** — every storage interface embeds `Streamable`, which
  exposes `Dump` and `Load` so any backend can be snapshotted and restored for
  the backup and restore features.
- **Test doubles** — `MockAuthStorage`, `MockConfigStorage` and `MockLogStorage`
  (constructed with `NewMockAuthStorage`, `NewMockConfigStorage` and
  `NewMockLogStorage`) are `testify` mocks consumers can use to exercise their
  logic without a real database.

## Layout

- **auth.go / config.go / log.go** — the three domain storage interfaces.
- **streamable.go** — the `Streamable` snapshot/restore contract embedded by
  every interface.
- **mock_\*.go** — the mock implementations used in tests.
- **[backends](backends)** — the concrete implementations of these interfaces.
