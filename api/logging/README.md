# logging

The package at `api/logging` provides the logging building blocks specific to
FlowG's REST API.

It exists so API code never reaches for the global logger directly. Routing all
API logs through a single helper guarantees they carry a common `api` channel
tag, which lets operators tell API activity apart from the rest of FlowG's
output when filtering logs.

The logger is resolved lazily on each call, so it always reflects the logging
configuration in force when a request is handled rather than the one present at
startup.

The request-scoped context helpers and the correlation-id aware `slog.Handler`
that this package builds on top of live in the common
[`internal/app/logging`](../../internal/app/logging) package.

## Contents

- **`main.go`** — `Logger`, the lazily-resolved logger tagged with the `api`
  channel.
