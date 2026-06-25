# mgmt

The package at `internal/services/mgmt` is FlowG's management server. It exposes
operational endpoints on a port separate from the main HTTP service, so they can
be firewalled independently of user traffic.

## Responsibilities

- **Liveness** — serves `GET /health` for health checks and orchestration
  probes.
- **Metrics** — exposes the Prometheus registry at `/metrics`.
- **Profiling** — mounts the `net/http/pprof` endpoints under `/debug/pprof/`,
  but only in builds compiled with the `debug` build tag.
- **Wiring** — `NewServer` returns an fx module binding the listener to the
  application lifecycle, optionally over TLS.

## Layout

- **main.go** — `ServerOptions`, the `Server`, and the `NewServer` fx module.
- **profiler_debug.go** — registers the pprof endpoints (build tag `debug`).
- **profiler_nodebug.go** — a no-op profiler registration for normal builds
  (build tag `!debug`).
