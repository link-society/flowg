# services

The packages under `internal/services` are FlowG's network-facing servers: the
long-lived listeners that accept connections from the outside world and hand the
work off to the engines and storage.

Each service is wired as an [fx](https://uber-go.github.io/fx/) module that binds
its listener (and, where relevant, its worker actor) to the application
lifecycle, so the top-level command can compose exactly the services it needs.

## Layout

- **[http](http)** — the main HTTP server; mounts the REST API and the web UI
  behind a single port and the access-log middleware.
- **[syslog](syslog)** — a syslog listener (UDP/TCP, optionally TLS) that feeds
  received messages into the pipeline engine.
- **[mgmt](mgmt)** — the operational server exposing health, Prometheus metrics
  and (in debug builds) the pprof profiler on a separate port.
