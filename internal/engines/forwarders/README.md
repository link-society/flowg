# forwarders

The package at `internal/engines/forwarders` executes FlowG's forwarders: it
turns a stored forwarder configuration (`models.ForwarderV2`) into a `Runtime`
that delivers log records to the configured external destination.

It exists to keep behaviour out of the data model. The `internal/models`
package describes _what_ a forwarder is (its JSON shape, validation and schema
hints); this package owns _how_ each backend is contacted, along with the
backend SDKs that entails. Consumers that only need the shapes — the API
schemas, the CLI client — no longer pull in every destination's client library.

## The Runtime interface

`NewRuntime` inspects the configuration's tagged union and returns the matching
implementation (or `ErrNotImplemented`). Every runtime follows the same
lifecycle:

- **Init** — compiles the configuration's dynamic fields and builds the backend
  client.
- **Call** — delivers one record to the destination.
- **Close** — releases the connection, for the backends that hold one.

The pipelines engine builds a runtime for each forward node when it compiles a
flow, and calls it once per record; the test-forwarder API endpoint does the
same with a sample record.

## Dynamic fields

Several backends let parts of the payload (message, tags, routing key, ...) be
computed per record. `CompileDynamicField` turns such a value into an
[expr](https://expr-lang.org/) program: a value prefixed with `@expr:` is
compiled as an expression, anything else as a literal that evaluates to itself.
At call time the program runs against the record, exposed as `log` (the field
map) and `timestamp`.

## Layout

- **main.go** — the `Runtime` interface, the `NewRuntime` dispatcher and
  `CompileDynamicField`.
- one file per backend — **http**, **syslog**, **datadog**, **amqp**,
  **splunk**, **otlp**, **elastic**, **clickhouse**, **awscloudwatch** and
  **googlecloudlogging** — each implementing `Runtime` for the matching
  `models.Forwarder*V2` configuration.
