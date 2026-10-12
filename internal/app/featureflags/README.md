# featureflags

The package at `internal/app/featureflags` defines `Flags`, FlowG's feature
toggles. A `Flags` value is provided to the dependency graph from the
configuration, and injected into the components that depend on a toggle.

## Flags

- **Demo mode** — `Flags.DemoMode`. When enabled, the bootstrap seeds a "demo"
  role and user, every stream is given a short retention (15 minutes, 10MB) and
  no index, and the web UI is told so through its feature flags.
