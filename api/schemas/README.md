# schemas

The package at `api/schemas` declares the request and response types of FlowG's
REST API — one file per endpoint, mirroring the [operations](../operations)
package.

It exists so that API clients can depend on the wire shapes alone. The
`flowg-client` binary builds its requests from these types; because they live
apart from the operations, importing them does not drag in the operations'
dependencies (storage backends, engines, fx wiring). The operations themselves
import this package to bind their interactors to the same shapes the clients
use.

## Conventions

Each endpoint contributes a pair of types:

- a request struct (`XRequest`) describing the inputs, with tags mapping fields
  to path, query, header or body parameters;
- a response struct (`XResponse`) describing the output, usually carrying a
  `Success` flag alongside the payload.

The struct tags double as OpenAPI metadata: the API layer generates the JSON
schema of every endpoint from these definitions, so the documentation always
matches what the server actually decodes. Domain objects embedded in the
payloads (users, roles, flow graphs, forwarders, ...) come from
`internal/models`.

## Special cases

Most types are plain data, but a few carry protocol behaviour:

- **ingest_logs_otlp.go** — `IngestLogsOTLPRequest` implements the usecase
  loader interface to decode the OTLP payload itself, handling both protobuf
  and JSON encodings and optional gzip compression.
- **watch_logs.go** — `WatchLogsResponse` embeds the response writer so the
  operation can emit a Server-Sent Events stream instead of a single buffered
  body.
- **backup_*.go / restore_*.go** — the backup responses and restore requests
  stream whole database snapshots rather than JSON documents.
