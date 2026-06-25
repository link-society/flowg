# flowg-health

The `flowg-health` binary is a self-contained health check for a running
`flowg-server` process. It is designed to be used as a container `HEALTHCHECK`.

Rather than being configured separately, it discovers how to reach the server by
reading the target process's command line: it inspects the `--mgmt-bind` and
`--mgmt-tls` flags of the `flowg-server` process (looked up by PID) so it always
probes the right management endpoint, then issues a request to `/health` and
exits non-zero if the server is unhealthy.

## Layout

- **`main.go`** — the entrypoint; it builds the root command and propagates the
  exit code.
- **[cmd](./cmd)** — the root `flowg-health` command: process inspection and the
  health probe.
