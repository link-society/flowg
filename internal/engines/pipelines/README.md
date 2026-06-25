# pipelines

The package at `internal/engines/pipelines` runs FlowG's pipelines. A pipeline
is authored in the UI as a flow graph; this engine compiles that graph into an
executable node graph and drives log records through it.

The engine is built on a single [actor](https://github.com/vladopajic/go-actor)
that owns a cache of compiled pipelines, so a pipeline is built from storage once
and reused until its definition changes.

## Responsibilities

- **Compilation** — turns a stored `FlowGraphV2` into a graph of `Node`s,
  resolving the transformers and forwarders each node references.
- **Execution** — feeds a record into a chosen entrypoint and propagates it
  through the graph, concurrently fanning out to each node's successors.
- **Caching** — keeps compiled pipelines hot and exposes invalidation so changes
  to a pipeline (or its dependencies) take effect on the next run.
- **Tracing** — supports dry runs that record what each node received, emitted
  and errored, without performing side effects.
- **Wiring** — `NewRunner` returns an `fx` module providing a `Runner` bound to
  the application lifecycle.

## Layout

- **main.go** — the `Runner` interface, its actor-backed implementation and `fx`
  wiring.
- **worker.go** — the actor body; owns the compiled-pipeline cache.
- **messages.go** — the actor messages (run, invalidate one, invalidate all) and
  the entrypoint constants.
- **types_pipeline.go** — `Pipeline` and the `BuildFlow`/`BuildFromStorage`
  compilers.
- **types_nodes.go** — the `Node` interface and the node implementations
  (source, transform, switch, pipeline, forward, router).
- **node_tracer.go** — the dry-run tracer carried through the context.
- **context.go** — context plumbing used to reach the worker from inside nodes.
- **errors.go** — the typed errors raised while compiling a flow graph.
- **mock.go** — a testify mock of `Runner` for tests.

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
pipeline) skip their effects and every node appends a `NodeTrace`, letting the UI
show exactly how a record would travel through the pipeline.
