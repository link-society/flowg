# config

The package at `internal/app/config` defines the user-facing configuration of
`flowg-server`: the schema of the HCL configuration file, its defaults, and its
validation.

It is deliberately self-contained: it does not import the rest of the
application (and therefore does not depend on the FoundationDB C library), so
the configuration can be loaded and tested on its own. Converting it into the
options of each component is the job of the caller.

## Layout

- **`schema.go`** — the schema (`RootConfig` and its sections), `DefaultConfig`,
  and the `Load` / `Validate` methods that parse and check a configuration file.
- **`env.go`** — default settings, each resolved from an environment variable
  via the `getEnv*` helpers.
