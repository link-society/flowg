# featureflags

The package at `internal/app/featureflags` holds FlowG's process-wide feature
toggles. The flags are stored in atomics, so they can be read and written from
any goroutine without extra synchronisation.

## Flags

- **Demo mode** — `GetDemoMode` / `SetDemoMode`. When enabled, FlowG disables
  state-mutating operations so a public instance can be exposed safely.
