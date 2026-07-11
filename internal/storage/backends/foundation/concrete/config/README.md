# config

The package at `internal/storage/backends/foundation/concrete/config` wires the
backend-agnostic configuration store from
[databases/config](../../../../databases/config) onto the FoundationDB
[kv.Adapter](../../../../generic/kv).

It assembles the FoundationDB adapter (scoped to the `config` subspace) with the
`ConfigStorage` implementation into a single `fx` module that provides the
`ConfigStorage` interface declared in [interfaces](../../../../interfaces).

## Wiring

`NewStorage` returns an `fx` module that provides a `ConfigStorage`. `Options`
carries the FoundationDB cluster file and the shared key space, and
`DefaultOptions` supplies their defaults.

## Layout

- **main.go** — the `Options`, `DefaultOptions` and `NewStorage` `fx` wiring.
