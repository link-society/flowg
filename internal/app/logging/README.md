# logging

The package at `internal/app/logging` holds the logging building blocks shared
between FlowG's REST API and the `flowg-server` process.

It exists so the request-scoped correlation handling and the `slog` handler that
both layers rely on live in one place: the API ([`api/logging`](../../../api/logging))
builds its access-log middleware on top of these helpers, while the server
([`cmd/flowg-server/logging`](../../../cmd/flowg-server/logging)) uses them to
configure the process-wide logger.

## Contents

- **`context.go`** — request-scoped context helpers: `WithCorrelationId` and
  `WithSensitiveMarker` enrich a request context, while `MarkSensitive` and
  `IsMarkedSensitive` let handlers flag and detect requests that carry sensitive
  data.
- **`handler.go`** — `NewHandler`, a `slog.Handler` that attaches the request
  `correlation_id` to every record emitted within the request scope, plus the
  `VERBOSE_LOGGING` flag that toggles dumping response bodies for failed
  requests.
