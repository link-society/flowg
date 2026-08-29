# pipelines

The package at `internal/engines/pipelines` runs FlowG's pipelines. A pipeline
is authored in the UI as a flow graph; this engine compiles that graph into an
executable node graph and drives log records through it.

The engine is built on a single [actor](https://github.com/vladopajic/go-actor)
that owns the set of running pipeline builds. Pipelines are eager: every
persisted pipeline is compiled and started at boot, and the
[confignotify](../confignotify) events emitted by the config storage rebuild or
terminate builds as the configuration changes, wherever the change was made.

## Responsibilities

- **Compilation** — turns a stored `FlowGraphV2` into a graph of `Node`s,
  resolving the transformers each node references and building a
  [forwarders](../forwarders) runtime for each forward node.
- **Execution** — feeds a record into a chosen entrypoint and propagates it
  through the graph, concurrently fanning out to each node's successors.
- **Lifecycle** — starts every persisted pipeline at boot (a broken pipeline is
  logged and skipped so the server stays up), swaps in new builds on
  `PipelineChanged`, retires builds on `PipelineDeleted`, and reconciles
  everything on `DependenciesChanged`.
- **Tracing** — supports dry runs that record what each node received, emitted
  and errored, without performing side effects.
- **Wiring** — `NewRunner` returns an `fx` module providing a `Runner` bound to
  the application lifecycle.

## Build lifetime

Every use of a running build (a record being processed, a metrics scrape, a
nested pipeline call) holds an in-flight reservation on it. Retiring a build —
on update, deletion or shutdown — first marks it closed so new reservations
fail, then waits for the in-flight ones to drain before closing node resources
(forwarder connections, metric collectors). On update, in-flight records finish
on the old build while new records already flow through the new one; if the new
build fails to compile, the old one keeps running. The drain wait is bounded by
the caller's context, but the teardown itself always completes in the
background.

## Layout

- **main.go** — the `Runner` interface (`Process`, `ScrapMetrics`), its
  actor-backed implementation, the internal run/terminate lifecycle and the
  `fx` wiring (eager start at boot, subscription to config events, graceful
  stop).
- **worker.go** — the actor body; owns the running builds, the acquire/retire
  lifetime logic and the event-driven reconcile.
- **messages.go** — the actor messages (process, scrap metrics, run, terminate)
  and the entrypoint constants.
- **types_pipeline.go** — `Pipeline` and the `BuildFlow`/`BuildFromStorage`
  compilers.
- **types_nodes.go** — the `Node` interface and the node implementations
  (source, transform, switch, pipeline, forward, router).
- **node_tracer.go** — the dry-run tracer carried through the context.
- **context.go** — context plumbing used to reach the worker from inside nodes.
- **errors.go** — the typed errors raised while compiling a flow graph.
- **mocks/** — a testify mock of `Runner` for tests.

## Node types

A compiled pipeline is a directed graph of nodes:

- **source** — an entrypoint; forwards records to its successors. Each source's
  declared type (e.g. `direct`, `syslog`) names an entrypoint.
- **transform** — runs a VRL transformer, which may emit zero, one or many
  records per input.
- **switch** — forwards a record only when it matches an
  [expr](https://expr-lang.org/) condition.
- **pipeline** — delegates processing to another named pipeline.
- **forward** — sends the record to an external destination through a forwarder.
- **router** — a terminal node; persists the record to a stream and notifies
  live subscribers.

## Dry runs

When a `NodeTracer` is attached to the context (via `WithTracer`), processing
switches to dry-run mode: side-effecting nodes (forward, router, nested
pipeline) skip their effects and every node appends a trace entry
(`models.PipelineNodeTrace`), letting the UI show exactly how a record would
travel through the pipeline.
