# logging

The package at `cmd/flowg-server/logging` configures the process-wide logger for
the `flowg-server` process.

It exists to keep global logging setup out of the API and engine packages: the
server decides, at startup, where logs go and how verbose they are, while the
rest of the codebase only ever reaches for `slog`'s default logger.

Both helpers install a handler from
[`internal/app/logging`](../../../internal/app/logging) so the correlation id
carried by each request is attached to every record.

## Contents

- **`main.go`** — `Setup`, which wires the default `slog` logger to standard
  output at the requested verbosity, and `Discard`, which silences all output
  (used by tests).
