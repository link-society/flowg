# middlewares

The package at `api/middlewares` exposes FlowG's log ingestion under the wire
protocols of popular log shippers.

It exists so existing agents can forward logs to FlowG without being
reconfigured to speak its native API. Each middleware translates a foreign
protocol into the same internal action — authenticate the caller, then run the
submitted logs through the addressed pipeline — keeping the compatibility shims
isolated from the rest of the API.

## Responsibilities

- **Elasticsearch compatibility** — accepts the subset of the Elasticsearch API
  that log shippers exercise: advertising itself as Elasticsearch so official
  clients agree to connect, checking that an index (a FlowG pipeline) exists,
  and indexing a document (running a log record through that pipeline).
- **Authentication** — resolves the HTTP Basic credentials these clients send
  into a FlowG user, so each request is subject to the same permission checks as
  the native API.

## Usage shape

Each middleware registers itself with the [routing](../routing) table from an
`init` function, binding its handler to the prefix it is mounted under. The
`api` package collects the whole set and mounts each one alongside its
operations, so importing this package is all it takes to enable the middlewares.

```text
init() ──▶ routing.RegisterMiddleware ──▶ Middleware{Pattern, Handler}

shipper ──▶ elasticAuth ──▶ permission check ──▶ pipeline run
            (authenticate)   (authorize)          (ingest)
```
