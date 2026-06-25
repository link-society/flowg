# bootstrap

The package at `internal/storage/bootstrap` seeds a fresh FlowG deployment with
the minimal data it needs to be usable.

It exists to keep first-run provisioning in one place, written against the
[storage](..) contracts rather than any concrete backend. The server runs these
helpers at startup so an empty database is brought to a known, working state
without operator intervention; each helper is idempotent and only writes when
the relevant data is missing, making it safe to run on every boot.

## Responsibilities

- **Default roles and users** — `DefaultRolesAndUsers` provisions the initial
  administrator role and account (and can reset a user's password) when the auth
  storage holds no roles yet.
- **Default system configuration** — `DefaultSystemConfig` writes the initial
  system configuration, including the allowed Syslog origins, when none exists.
- **Default pipeline** — `DefaultPipeline` creates a `default` pipeline so logs
  can be ingested out of the box.

## Layout

- **auth.go** — `DefaultRolesAndUsers` plus its `BootstrapAuthOptions` and
  `ResetUserOptions` inputs.
- **system.go** — `DefaultSystemConfig` and its `BootstrapSystemOptions` input.
- **pipelines.go** — `DefaultPipeline`.
