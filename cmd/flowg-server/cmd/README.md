# cmd

The package at `cmd/flowg-server/cmd` defines the root command of the
`flowg-server` binary.

It exists to keep the process's command-line surface in one place: the flag
definitions, the loading of the
[`internal/app/config`](../../../internal/app/config) configuration, and the
startup of the [`internal/app/server`](../../../internal/app/server) application
built from it.

## Layout

- **`main.go`** — `NewRootCommand`, which assembles the root command, loads and
  validates the configuration, runs the startup hooks (umask, logging,
  metrics) and starts the fx application.
