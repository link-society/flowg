# concrete

The packages under `internal/storage/backends/badger/concrete` wire FlowG's
domain storage onto BadgerDB.

Each package assembles the shared BadgerDB [kv.Adapter](../../../generic/kv) with
the matching backend-agnostic store from [databases](../../../databases) into a
single `fx` module that provides one of the interfaces declared in
[interfaces](../../../interfaces). They also run the BadgerDB-specific startup
migrations for their domain. Keeping them here separates that wiring from the
shared BadgerDB primitives that sit alongside them (the adapter, the transaction
adapter, and the `slog` logger).

Every package exposes the same shape: an `Options` struct, a `DefaultOptions`
constructor and a `NewStorage` function returning the `fx` module.

## Layout

- **[auth](auth)** — provides `AuthStorage`, and migrates the legacy "alerts"
  permission scopes to "forwarders" on startup.
- **[config](config)** — provides `ConfigStorage`, and migrates legacy on-disk
  configuration (including the "alerts" directory) into BadgerDB on startup.
- **[log](log)** — provides `LogStorage`, and starts the retention garbage
  collector.
