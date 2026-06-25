# cmd

The package at `cmd/flowg-health/cmd` defines the root command of the
`flowg-health` binary.

It exists to keep the health-check logic in one place: locating the target
`flowg-server` process, deriving its management endpoint from the process's own
flags, and probing `/health`.

## Layout

- **`main.go`** — `NewRootCommand`, which reads the target process's command
  line (by `--pid`), resolves the management address and TLS mode, and performs
  the HTTP health probe.
- **`env.go`** — the management-address defaults, resolved from environment
  variables via the `getEnv*` helpers (kept in sync with `flowg-server`).
