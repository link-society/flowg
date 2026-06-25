# http

The package at `internal/services/http` is FlowG's main HTTP server. It binds a
single port and serves both the REST API and the web UI from it, behind a shared
access-log middleware.

## Responsibilities

- **Composition** — mounts the [api](../../../api) handler under `<mount>/api/`
  and the [web](../../../web) handler under `<mount>/web/`, redirecting the root
  to the UI.
- **Transport** — listens on the configured address, optionally over TLS, and is
  started and stopped with the application lifecycle.
- **Access logging** — wraps every request with a correlation id and emits a
  structured access-log record (see below).

## Layout

- **main.go** — `ServerOptions`, the `Server`, and the `NewServer` fx module that
  assembles the routes and lifecycle hooks.
- **logging.go** — the access-log middleware.

## Access logging

`loggingMiddleware` assigns each request a correlation id (from the
`X-Correlation-Id` header, or generated when absent), propagates it through the
request context, and logs one structured record on the `accesslog` channel once
the handler returns. Requests marked sensitive are logged at debug level rather
than info. When verbose logging is enabled, the response bodies of failed
requests (status >= 400) are buffered and dumped to standard error, framed by
their correlation id.
