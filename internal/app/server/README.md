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
- **Configuration** — `Options` gathers the bind addresses, TLS settings,
  storage directories and initial/reset credentials in one place.
- **Bootstrap** — registers a handler that, on start, seeds the default system
  configuration, roles and users, and pipeline, and applies the optional
  admin-credential reset (see [storage/bootstrap](../../storage/bootstrap)).

## Layout

- **main.go** — `Options` and the `NewServer` fx module.
- **bootstrap.go** — the `bootstrapHandler` run from the module's `OnStart` hook.
