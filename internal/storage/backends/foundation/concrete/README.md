# concrete

The packages under `internal/storage/backends/foundation/concrete` wire FlowG's
domain storage onto FoundationDB.

Each package assembles the shared FoundationDB
[kv.Adapter](../../../generic/kv) with the matching backend-agnostic store from
[databases](../../../databases) into a single `fx` module that provides one of
the interfaces declared in [interfaces](../../../interfaces). Keeping them here
separates that wiring from the shared FoundationDB primitives that sit alongside
them (the adapter, the transaction adapters and the garbage collector).

Every package exposes the same shape: an `Options` struct carrying the
FoundationDB cluster file and the shared key space, a `DefaultOptions`
constructor and a `NewStorage` function returning the `fx` module. Each storage
is scoped to its own subspace (`auth`, `config` or `log`) beneath the key space.

## Layout

- **[auth](auth)** — provides `AuthStorage`.
- **[config](config)** — provides `ConfigStorage`.
- **[log](log)** — provides `LogStorage`, and starts the retention garbage
  collector.
