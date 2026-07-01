# cmd

The package at `cmd/flowg-server/cmd` defines the root command of the
`flowg-server` binary.

It exists to keep the process's command-line surface in one place: the flag
definitions, their environment-variable defaults, and the logic that turns the
collected options into the [`internal/app/server`](../../../internal/app/server)
configuration the application is built from.

## Layout

- **`main.go`** — `NewRootCommand`, which assembles the root command, runs the
  startup hooks (umask, logging, demo mode, metrics) and starts the fx
  application.
- **`env.go`** — the per-flag defaults, each resolved from an environment
  variable via the `getEnv*` helpers.
- **`config.go`** — parses the configuration file
