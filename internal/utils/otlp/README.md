# otlp

The package at `internal/utils/otlp` decodes OpenTelemetry (OTLP) logs export
requests into FlowG log records.

It exists to keep the OTLP wire format out of the ingestion handlers. The
handlers hand it a request body and receive ready-to-store
`models.LogRecord` values, so the protobuf schema and its JSON encoding are
dealt with in a single place.

## Responsibilities

- **JSON ingestion** — `UnmarshalJSON` decodes an OTLP/HTTP logs request encoded
  as protobuf JSON.
- **Protobuf ingestion** — `UnmarshalProtobuf` decodes the same request encoded
  as binary protobuf.
- **Conversion** — both entry points flatten the resource/scope/record hierarchy
  of an `ExportLogsServiceRequest` into a flat slice of FlowG log records,
  promoting the well-known OTLP fields to top-level names and prefixing each
  attribute with `attr.`.

## Scope

Only the OTLP **logs** signal is supported; metrics and traces are out of scope.
