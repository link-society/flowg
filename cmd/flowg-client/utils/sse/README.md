# sse

The package at `cmd/flowg-client/utils/sse` reads
[Server-Sent Events](https://html.spec.whatwg.org/multipage/server-sent-events.html)
from a byte stream.

It exists for the streaming commands of the client (`stream tail`,
`stream watch`, `stream history`), which consume FlowG's event stream of log
records. It parses the wire format — events separated by a blank line, each a
set of `field: value` lines — into `Event` values, leaving the interpretation of
the payload to the caller.

## Vocabulary

- **`Event`** — a single decoded event: its `ID`, its `Type`, and its raw
  `Data`.
- **`EventStreamReader`** — reads from the stream one event at a time; `Next`
  returns the next `Event` or `io.EOF` once the stream ends.
