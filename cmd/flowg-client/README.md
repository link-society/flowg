# flowg-client

The `flowg-client` binary is a command-line client for FlowG's HTTP API.

It exists to drive a running FlowG instance from a shell or a script — managing
streams, pipelines, transformers, forwarders, access control and tokens, and
tailing logs in real time — without going through the web UI.

## Layout

- **`main.go`** — the entrypoint; it builds the root command and propagates the
  exit code chosen by whichever subcommand ran.
- **[cmd](./cmd)** — the command tree itself, one file per subcommand, wired
  together under the root `flowg-client` command.
- **[utils](./utils)** — the building blocks the commands share: the
  authenticated HTTP client, custom flag types, the log printer, and the
  Server-Sent Events reader used by the streaming commands.

## Configuration

The API endpoints and the authentication token are read from flags, each
defaulting to an environment variable so the client can be configured once for a
session:

- `--api` / `FLOWG_API` — the FlowG HTTP API.
- `--api-token` / `FLOWG_API_TOKEN` — the bearer token used to authenticate.
- `--mgmt-api` / `FLOWG_MGMT_API` — the management HTTP API (used for backups).
