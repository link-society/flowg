# cmd

The package at `cmd/flowg-server/cmd` defines the root command of the
`flowg-server` binary.

It exists to keep the process's command-line surface in one place: the flag
definitions, and the logic that turns the loaded
[`internal/app/config`](../../../internal/app/config) configuration into the
[`internal/app/server`](../../../internal/app/server) options the application is
built from.

## Layout

- **`main.go`** — `NewRootCommand`, which assembles the root command, loads and
  validates the configuration, runs the startup hooks (umask, logging, demo
  mode, metrics) and starts the fx application.
- **`server_options.go`** — converts the configuration into the
  `server.Options` struct, loading the TLS certificates of the services that
  enable TLS.
