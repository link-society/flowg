# flowg-server

The `flowg-server` binary is the FlowG server process.

It is the long-running daemon that exposes the platform: the HTTP API and web UI,
the syslog ingestion listener, and the management endpoints. It reads its
configuration from flags (each backed by an environment variable), wires the
whole application together and runs it until interrupted.

## Layout

- **`main.go`** — the entrypoint; it builds the root command and propagates the
  exit code chosen by the command.
- **[cmd](./cmd)** — the root `flowg-server` command: flag definitions,
  environment-variable defaults, and the translation from CLI options into the
  application's configuration.
- **[logging](./logging)** — process-wide `slog` logger setup, sharing the
  correlation-id handler with the API.

## Configuration

Every flag defaults to an environment variable so the server can be configured
entirely through the environment. The main groups are:

- **HTTP** (`--http-*` / `FLOWG_HTTP_*`) — bind address, mount path and TLS for
  the API and web UI.
- **Management** (`--mgmt-*` / `FLOWG_MGMT_*`) — bind address and TLS for the
  health/metrics server.
- **Syslog** (`--syslog-*` / `FLOWG_SYSLOG_*`) — protocol, bind address, TLS and
  initial allowed origins for the syslog listener.
- **Storage** (`--auth-dir`, `--config-dir`, `--log-dir` / `FLOWG_*_DIR`) — the
  three on-disk database directories.
- **Bootstrap** (`--auth-initial-*`, `--auth-reset-*`) — the initial admin
  account and the optional password reset applied on start.
