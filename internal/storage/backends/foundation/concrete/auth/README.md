# auth

The package at `internal/storage/backends/foundation/concrete/auth` wires the
backend-agnostic authentication store from
[databases/auth](../../../../databases/auth) onto the FoundationDB
[kv.Adapter](../../../../generic/kv).

It assembles the FoundationDB adapter (scoped to the `auth` subspace) with the
`AuthStorage` implementation into a single `fx` module that provides the
`AuthStorage` interface declared in [interfaces](../../../../interfaces).

## Wiring

`NewStorage` returns an `fx` module that provides an `AuthStorage`. `Options`
carries the FoundationDB cluster file and the shared key space, and
`DefaultOptions` supplies their defaults.

## Layout

- **main.go** — the `Options`, `DefaultOptions` and `NewStorage` `fx` wiring.
