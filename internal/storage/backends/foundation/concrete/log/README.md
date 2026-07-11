# log

The package at `internal/storage/backends/foundation/concrete/log` wires the
backend-agnostic log store from [databases/log](../../../../databases/log) onto
the FoundationDB [kv.Adapter](../../../../generic/kv).

It assembles the FoundationDB adapter (scoped to the `log` subspace) with the
`LogStorage` implementation into a single `fx` module that provides the
`LogStorage` interface declared in [interfaces](../../../../interfaces), and
starts the retention garbage collector.

## Wiring

`NewStorage` returns an `fx` module that provides a `LogStorage` and starts a
background worker enforcing each stream's retention budget. `Options` carries
the FoundationDB cluster file, the shared key space and the garbage-collection
interval, and `DefaultOptions` supplies their defaults.

Two collectors run for this storage: the retention collector started here, and
the adapter's own expired-key sweeper (FoundationDB has no native TTL). Both use
the configured `GCInterval`.

## Layout

- **main.go** — the `Options`, `DefaultOptions` and `NewStorage` `fx` wiring.
