# utils

The package at `cmd/flowg-client/utils` holds the building blocks the
`flowg-client` commands share.

It exists to keep the command files focused on their own logic: the concerns
that every command needs — talking to the API, parsing repeated flags, and
printing log records — live here once instead of being repeated across the
command tree.

## Contents

- **`http.go`** — `Client`, a minimal HTTP client that targets a FlowG API base
  URL and authenticates each request with a bearer token.
- **`flags.go`** — `IndexMap`, a flag type that collects repeated `key=value`
  arguments into a multimap.
- **`printer.go`** — `Printer`, which renders log records to standard output as
  logfmt lines with fields in a stable order.
- **[sse](./sse)** — the Server-Sent Events reader used by the streaming
  commands (`stream tail`, `stream watch`, `stream history`).
