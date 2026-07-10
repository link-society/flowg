# storage

The package tree under `internal/storage` defines FlowG's storage contracts and
the implementations that satisfy them.

It exists to decouple the application from any particular database technology.
The engines, services and API operations are written against the interfaces
declared here; the concrete backends that satisfy them are wired in at startup.
Swapping or adding a backend therefore never requires touching the consumers.

## Responsibilities

- **Domain contracts** — `AuthStorage`, `ConfigStorage` and `LogStorage`
  describe the operations FlowG needs to manage, respectively, identities and
  permissions, pipeline/transformer/forwarder configuration, and ingested log
  records.
- **Shared capability** — every storage interface embeds `Streamable`, which
  exposes `Dump` and `Load` so any backend can be snapshotted and restored for
  the backup and restore features.
- **Reusable logic** — the domain stores are implemented once against a generic
  key-value abstraction, so a backend only has to provide that abstraction rather
  than reimplement the domain logic.
- **First-run provisioning** — bootstrap helpers seed an empty deployment with
  the default roles, users, system configuration and pipeline it needs to be
  usable.
- **Test doubles** — mock implementations of the three interfaces let consumers
  exercise their logic without a real database.

## Layout

- **[interfaces](interfaces)** — the domain storage contracts (`AuthStorage`,
  `ConfigStorage`, `LogStorage`) and the `Streamable` snapshot/restore contract
  they embed.
- **[generic](generic)** — backend-agnostic building blocks, chiefly the
  [kv](generic/kv) key-value store abstraction the domain stores are written
  against.
- **[databases](databases)** — backend-agnostic implementations of the domain
  interfaces, expressed in terms of the generic key-value abstraction.
- **[backends](backends)** — the concrete backends that instantiate the domain
  stores with a real database adapter.
- **[bootstrap](bootstrap)** — idempotent helpers that seed a fresh deployment
  with its default data.
- **[mocks](mocks)** — `testify`-based mock implementations of the three storage
  interfaces for use in tests.
