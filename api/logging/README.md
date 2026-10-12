# logging

The package at `api/logging` provides the logging building blocks specific to
FlowG's REST API.

It exists so API code never reaches for a global logger. Every operation is
injected the root logger, and routes its logs through a single helper deriving
the API logger from it, which guarantees they carry a common `api` channel tag
and lets operators tell API activity apart from the rest of FlowG's output when
filtering logs.

The request-scoped context helpers and the correlation-id aware `slog.Handler`
that this package builds on top of live in the common
[`internal/app/logging`](../../internal/app/logging) package.

## Contents

- **`main.go`** — `Logger`, deriving the logger tagged with the `api` channel
  from the root logger.
