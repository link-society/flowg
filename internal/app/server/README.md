# server

The package at `internal/app/server` assembles the complete FlowG server. Its
`NewServer` function returns a single [fx](https://uber-go.github.io/fx/) module
that wires every layer together and binds them to the application lifecycle.

## Responsibilities

- **Composition** — instantiates and connects the three layers:
  - **Storage** — the auth, config and log [storage](../../storage) backends.
  - **Engines** — the [log notifier](../../engines/lognotify) and the
    [pipeline runner](../../engines/pipelines).
  - **Services** — the [http](../../services/http), [mgmt](../../services/mgmt)
    and [syslog](../../services/syslog) servers.
- **Configuration** — supplies the [configuration](../config) to the graph,
  selects the storage backend modules, and converts the configuration into the
  options of each component (see `wiring.go`).
- **Bootstrap** — registers a handler that, on start, seeds the default system
  configuration, roles and users, and pipeline, and applies the optional
  admin-credential reset (see [storage/bootstrap](../../storage/bootstrap)).

## Layout

- **main.go** — the `NewServer` fx module.
- **wiring.go** — the configuration module (supplies the configuration and
  provides the options of each component from it) and the selection of the
  storage backend modules. It is the only place knowing both the configuration
  schema and the components' options.
- **bootstrap.go** — the `bootstrapHandler` run from the module's `OnStart`
  hook, configured by the `bootstrapOptions` provided to the graph.
