# config

The package at `internal/app/config` defines the user-facing configuration of
`flowg-server`: the schema of the HCL configuration file, its defaults, and its
validation.

It is deliberately self-contained: it does not import the rest of the
application (and therefore does not depend on the FoundationDB C library), so
the configuration can be loaded and tested on its own. Converting it into the
options of each component is the job of the caller.

## Layout

- **`schema.go`** — the schema (`RootConfig` and its sections), the built-in
  defaults (`DefaultConfig`), and the `Load` / `Validate` methods that complete
  the configuration from the configuration file and the environment, then check
  it.
- **`env.go`** — the environment overlay applied by `Load` after the
  configuration file (environment variables take precedence), and the `getEnv*`
  helpers reading a typed value from an environment variable, with a fallback
  when it is unset or empty.
