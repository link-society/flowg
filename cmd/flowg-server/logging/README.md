# logging

The package at `cmd/flowg-server/logging` builds the root logger of the
`flowg-server` process.

It exists to keep logging setup out of the API and engine packages: the server
decides, at startup, where logs go and how verbose they are, and injects the
resulting root logger into the application, from which every component derives
its own (tagged with its `channel`).

Both helpers install a handler from
[`internal/app/logging`](../../../internal/app/logging) so the correlation id
carried by each request is attached to every record.

## Contents

- **`main.go`** — `Setup`, which builds the root logger writing to standard
  output at the requested verbosity (and makes it `slog`'s default, for
  third-party code), and `Discard`, which silences `slog`'s default logger
  (used by tests).
