# metrics

The package at `internal/app/metrics` holds FlowG's Prometheus instrumentation.
It defines the process-wide counters and the helpers the engines call to record
activity, which the [management server](../../services/mgmt) exposes on its
`/metrics` endpoint.

## Responsibilities

- **Registration** — `Setup` creates the counters and registers them with the
  default Prometheus registry; it must be called once during startup.
- **Recording** — `IncStreamLogCounter` and `IncPipelineLogCounter` increment
  the per-stream ingestion counter and the per-pipeline processing counter
  (labelled `success`/`error`) respectively.

## Metrics

- `flowg_stream_log_total{stream}` — log records ingested per stream.
- `flowg_pipeline_log_total{pipeline,status}` — log records processed per
  pipeline, by outcome.
