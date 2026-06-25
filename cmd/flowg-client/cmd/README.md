# cmd

The package at `cmd/flowg-client/cmd` defines the command tree of the
`flowg-client` binary.

It exists to keep each subcommand self-contained: every verb the client exposes
lives in its own file and contributes a `cobra.Command`, so a command can be
added or changed without touching the others. The root command assembles them
and shares the configured API clients with each one.

## Layout

- **`main.go`** — builds the root command, declares the global flags (the API
  endpoints and token), and constructs the API clients once before any
  subcommand runs.
- **`context.go`** — the typed keys under which the configured API clients are
  stored on the command context, so every subcommand reaches them the same way.
- **`env.go`** — resolves each global flag's default from its environment
  variable.
- **`cmd_*.go`** — one file per command (and subcommand), grouped by resource:
  `acl`, `forwarder`, `pipeline`, `stream`, `transformer`, `token`,
  `systemconfig`, `admin`, and `login`.

## Usage shape

The root command builds the clients in its pre-run hook and stores them on the
context; each subcommand reads the client back out and calls the API.

```text
NewRootCommand ──▶ PersistentPreRun ──▶ context[ApiClient] = utils.NewClient(…)
                                                 │
                            cmd_*.go ──▶ ctx.Value(ApiClient) ──▶ client.Do(req)
```
