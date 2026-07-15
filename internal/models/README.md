# models

The package at `internal/models` holds FlowG's domain types: the plain data
structures that every other layer — storage, engines, services and the API —
shares to describe logs, pipelines, forwarders and authorization.

It is deliberately dependency-light and sits near the bottom of the import
graph. Types here own their JSON encoding, their OpenAPI/JSON-schema hints (via
struct tags and `JSONSchema*` hooks) and, where a format has evolved, the
conversion from older on-disk versions to the current one.

## Responsibilities

- **Canonical shapes** — define the records and configuration objects passed
  around the system, independent of how they are stored or transported.
- **Schema generation** — carry the tags and hooks that drive OpenAPI/JSON
  schema generation for the API.
- **Versioned migrations** — read legacy persisted formats and upgrade them to
  the current version on load.
- **Self-contained behaviour** — provide the small helpers that belong with the
  data (permission projection, scope parsing, version upgrades, ...). Behaviour
  that reaches outside the data — executing a forwarder, flattening OTLP
  payloads, building storage keys — lives with the engine or layer that owns it.

## Layout

### Authorization

- **auth_scope.go** — `Scope`, the atomic permissions, with parsing and
  enumeration.
- **auth_role.go** — `Role`, a named set of scopes.
- **auth_user.go** — `User`, an account with assigned roles.
- **auth_permissions.go** — `Permissions`, the boolean UI projection of a set of
  scopes.
- **auth_provider.go** — `AuthProvider` and the `AuthProviderConfig` tagged
  union for delegated authentication, with one file per backend
  (**auth_provider_oidc.go**, **auth_provider_saml.go**).

### Logs and streams

- **logrecord.go** — `LogRecord`, the canonical log entry: a timestamp plus a
  flat map of string fields.
- **stream_config.go** — `StreamConfig`, a stream's retention and indexing
  policy.
- **system_configuration.go** — `SystemConfiguration`, global server settings.

### Pipelines (flow graphs)

- **flow_v2.go** — `FlowGraphV2`, the current flow-graph shape compiled by the
  pipelines engine.
- **flow_v1.go** — the legacy V1 shape, kept for migration only.
- **flow_convert.go** — `ConvertFlowGraph`, which loads any supported version and
  returns the latest, upgrading switch conditions to expr-lang along the way.
- **pipeline_node_trace.go** — `PipelineNodeTrace`, the record of one node's
  execution during a pipeline dry run.

### Forwarders

- **forwarder_v2.go** — `ForwarderV2` and the `ForwarderConfigV2` tagged union
  that dispatches to one backend.
- **forwarder_v1.go** / **forwarder_convert.go** — the legacy V1 shape and its
  upgrade to V2.
- **forwarder_v2_*.go** — one file per backend (http, syslog, datadog, amqp,
  splunk, otlp, elastic, clickhouse, awscloudwatch, googlecloudlogging), each
  declaring its configuration shape and per-record field types (either a
  literal value or a `DynamicField`). Their execution lives in the
  [forwarders engine](../engines/forwarders).

### Helpers

- **dynamic_field.go** — `DynamicField`, a forwarder value that may be a literal
  or an [expr](https://expr-lang.org/) expression, compiled and evaluated per
  record by the forwarders engine.
