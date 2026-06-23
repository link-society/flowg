# routing

The package at `api/routing` wires API operations to their HTTP routes.

It exists to let an API package assemble its complete routing table out of
endpoints that each live in their own file, with no central list to keep in
sync. The route metadata travels with the behaviour, so an endpoint can be
added, moved or removed by editing a single file.

## Vocabulary

- **`Operation`** — the value that pairs an interactor with the route that
  reaches it, plus how it is exposed (public or authenticated, and any OpenAPI
  refinements the reflector cannot infer on its own).
- **`Middleware`** — the value that pairs a sub-handler with the route prefix it
  is mounted under, for protocol-compatibility shims that sit beside the native
  operations.
- **`RegisterOperation` / `RegisterMiddleware`** — how an endpoint or middleware
  contributes itself to the routing table, binding its constructor to a route.
  Both are called from an `init` function, so importing the package is enough to
  register them.
- **`Module`** — exposes every registered operation and middleware to the
  dependency-injection container as value groups, sparing callers from
  enumerating them.

## Usage shape

Each endpoint registers itself; the consuming API package collects the whole
group and mounts each operation on its route.

```text
init() ──▶ RegisterOperation ──▶ Operation{Method, Pattern, …} ──▶ group
                                                                     │
                                                  Module() ──▶ api package
```
