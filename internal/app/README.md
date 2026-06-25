# app

The package at `internal/app` is FlowG's application layer: the code that
composes storage, engines and services into a runnable process. It sits between
the `cmd/` entrypoints and the lower-level building blocks.

The `app` package itself only carries the generated `FLOWG_VERSION` constant
(produced from `VERSION.txt` via `go generate`); the actual application is built
out of its subpackages.

## Layout

- **[server](server)** — the `NewServer` fx module that wires the whole server
  together (storage, engines, services) and the bootstrap handler that seeds a
  fresh instance on start.
- **[logging](logging)** — request-scoped correlation helpers and the `slog`
  handler shared by the API and the server process.
- **[metrics](metrics)** — the Prometheus counters and helpers exposed on the
  management server's `/metrics` endpoint.
- **[featureflags](featureflags)** — process-wide feature toggles (e.g. demo
  mode).
