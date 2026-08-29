# config

The package at `internal/storage/backends/foundation/concrete/config` wires the
backend-agnostic configuration store from
[databases/config](../../../../databases/config) onto the FoundationDB
[kv.Adapter](../../../../generic/kv).

It assembles the FoundationDB adapter (scoped to the `config` subspace) with the
`ConfigStorage` implementation into a single `fx` module that provides the
`ConfigStorage` interface declared in [interfaces](../../../../interfaces).

## Wiring

`NewStorage` returns an `fx` module that provides a `ConfigStorage`. `Options`
carries the FoundationDB cluster file and the shared key space, and
`DefaultOptions` supplies their defaults.

The adapter is created with the change log enabled, and when a
[confignotify](../../../../../engines/confignotify) notifier is present in the
container the module also starts a watcher actor. The store itself emits no
local notification: every node — including the one that made the change —
observes mutations through the watcher, so local and remote changes follow the
same path.

## Change watcher

The watcher arms a FoundationDB watch on the adapter's change log key, hashes
the stored pipelines, transformers, forwarders and system configuration, and
diffs the result against its previous snapshot whenever the watch fires:
pipeline diffs become `PipelineChanged` / `PipelineDeleted` events, transformer
or forwarder diffs collapse into a single `DependenciesChanged`, and a system
configuration diff drops the node's cached system configuration. Arming the
watch before scanning guarantees no missed window, and diffing makes
FoundationDB's watch coalescing harmless. The first scan primes the snapshot
silently, since the pipeline runner builds everything at boot on its own.

## Layout

- **main.go** — the `Options`, `DefaultOptions` and `NewStorage` `fx` wiring.
- **watcher.go** — the change log watcher actor translating storage diffs into
  confignotify events.
