# engines

The packages under `internal/engines` are FlowG's runtime processing engines:
the long-lived, concurrent components that turn ingested log records into
transformed, routed and stored data.

Both engines are built on the [actor](https://github.com/vladopajic/go-actor)
model — each owns a single goroutine that serialises all of its state mutations —
and are wired into the application through [fx](https://uber-go.github.io/fx/)
modules whose lifecycle is bound to the server's.

## Layout

- **[pipelines](pipelines)** — compiles flow graphs into executable node graphs
  and runs records through them (transform, switch, forward, route).
- **[lognotify](lognotify)** — a live fan-out bus that pushes newly ingested
  records to interested subscribers, powering the live tail in the web UI.
