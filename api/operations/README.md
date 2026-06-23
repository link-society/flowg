# operations

The package at `api/operations` implements the operations of FlowG's REST API —
one per HTTP endpoint exposed through OpenAPI.

It exists to keep the behaviour of each endpoint separate from the HTTP wiring.
The surrounding `api` package owns routing, the OpenAPI service and the security
middlewares; this package owns _what_ each endpoint does once a request reaches
it. Splitting the two lets the wiring stay a flat, readable assembly while the
endpoints group naturally by domain.

## Anatomy of an operation

Each operation lives in its own file, declared in a fixed order:

- a dependency struct (`XDeps`), an [fx.In] struct listing only the storage
  backends and engines that operation actually uses;
- a request struct (`XRequest`) describing the inputs, with tags mapping fields
  to path, query, header or body parameters;
- a response struct (`XResponse`) describing the output;
- a constructor (`NewXUsecase`) that builds the interactor wiring inputs to
  outputs;
- an `init()` that calls `routing.RegisterOperation`, binding the constructor to
  its HTTP method, route pattern and any OpenAPI tweaks.

The routing primitives live in the [routing] package: `RegisterOperation`
records each endpoint as an `Operation` — the interactor bundled with its method,
pattern and metadata — and `Module` exposes the whole set as a dependency-
injection group. The `api` package collects that group and mounts each
`Operation` on its route. Because the route metadata travels with the behaviour,
an endpoint can be added, moved or removed by editing a single file.

```text
fx ──▶ NewXUsecase(XDeps) ──▶ interactor(XRequest) ──▶ XResponse
 (injects deps)  (per endpoint)   (business logic)

init() ──▶ routing.RegisterOperation ──▶ Operation{Method, Pattern, …} ──▶ group
```

[routing]: ../routing

## Authentication and authorization

Operations do not authenticate callers themselves — that is the job of the
`api/auth` middleware applied by the HTTP layer. Instead, each operation that
requires a permission wraps its interactor with a scope decorator from
`api/auth`, declaring the permission it needs. A handful of operations
(`login`, `whoami`, `change_password`, and personal-access-token management)
require only authentication, not a specific permission, and `login` requires
neither.

## Domains

The operations group by the resource they act on:

- **transformers** — read, list, save, delete and test VRL transformers.
- **pipelines** — read, list, save, delete and test flow-graph pipelines.
- **streams** — configure, inspect (fields, indices, usage), query, watch and
  purge log streams.
- **forwarders** — read, list, save, delete and test log forwarders.
- **ACLs** — manage roles, users and personal access tokens.
- **auth** — login, current-profile lookup and password change.
- **log ingestion** — push structured, textual or OpenTelemetry logs through a
  pipeline.
- **backup** — export and restore the authentication, configuration and log
  databases.
- **system configuration** — read and update the global system configuration.
