# api

The package at `api` assembles FlowG's REST API into a single HTTP handler.

It exists to keep endpoint behaviour and HTTP wiring apart. The behaviour and
route of every endpoint live in sibling packages; this package owns the concerns
that surround them — the OpenAPI service and its documentation, the routing
table that mounts each operation on its route, and the security middlewares that
authenticate callers — so assembly stays a flat, readable whole rather than
being scattered across endpoints.

## Composition

- **[operations](operations)** — one use-case per endpoint, each registering
  itself with the routing table.
- **[middlewares](middlewares)** — protocol-compatibility shims mounted beside
  the native operations.
- **[routing](routing)** — the vocabulary endpoints use to contribute
  themselves, collected here into the route table.
- **[auth](auth)** — the authentication middleware guarding every non-public
  route, and the authorization decorators endpoints apply.
- **[logging](logging)** — the shared logger every endpoint emits through.

Operations are mounted in a deterministic order so the generated OpenAPI
document is reproducible: the collector derives each shared component schema
from the first operation that references it, so the order must not depend on how
the container happens to yield the group.
