# langs

The packages under `internal/utils/langs` provide the small embedded languages
FlowG exposes to its users.

They exist to give pipelines and queries their expressiveness without coupling
that to any single engine. Each sub-package wraps a third-party language behind
a narrow, FlowG-specific API, so the rest of the application works with compiled
programs rather than raw expression strings.

## Layout

- **[filtering](filtering)** — the boolean expression language used to filter
  log records when querying a stream.
- **[vrl](vrl)** — a binding to Vector Remap Language, used by transformer nodes
  to reshape log events during pipeline processing.
