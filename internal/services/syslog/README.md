# syslog

The package at `internal/services/syslog` is a syslog ingestion server. It
listens for syslog messages and feeds each one into FlowG's pipeline engine,
giving syslog-emitting systems a direct path into the platform.

It is built on a worker [actor](https://github.com/vladopajic/go-actor) that
drains the messages parsed by the underlying syslog library.

## Responsibilities

- **Listening** — accepts messages over UDP, TCP or TCP+TLS, auto-detecting the
  syslog format (RFC 3164/5424).
- **Origin filtering** — when `SyslogAllowedOrigins` is configured, drops
  messages whose client IP is not an allowed address or within an allowed CIDR
  range.
- **Dispatch** — normalises each message into a `LogRecord` and runs it through
  every pipeline's `syslog` entrypoint concurrently.
- **Wiring** — `NewServer` returns an fx module binding the listener and actor to
  the application lifecycle.

## Layout

- **main.go** — `ServerOptions`, the `Server`, and the `NewServer` fx module that
  wires the listener, the message channel and the worker actor.
- **worker.go** — the actor body: origin filtering and per-pipeline dispatch.
- **parse.go** — converts the library's loosely-typed fields into a `LogRecord`.
