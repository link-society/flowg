# confignotify

The package at `internal/engines/confignotify` is the notification bus for
configuration changes. The config storage publishes an event after every
successful pipeline, transformer or forwarder mutation, and the pipeline runner
subscribes to (re)build or terminate the affected pipelines.

It exists to decouple the pipeline lifecycle from the API layer: mutations
trigger lifecycle changes at the storage level, so the runner reacts the same
way whether a change came from the local API, a bootstrap write, or (on
clustered backends) another node observed through a storage watcher.

## Responsibilities

- **Subscriptions** — registers per-subscriber mailboxes and tears them down
  when the subscriber's context is cancelled.
- **Fan-out** — broadcasts each event to every current subscriber, logging —
  but never failing on — individual delivery errors.
- **Wiring** — `NewNotifier` returns an `fx` module providing a `Notifier`
  whose actor and mailboxes follow the application lifecycle.

## Events

- **PipelineChanged** — a pipeline was created or updated; subscribers should
  rebuild it.
- **PipelineDeleted** — a pipeline was removed; subscribers should terminate it.
- **DependenciesChanged** — a transformer or forwarder changed, or the whole
  configuration was restored; subscribers should reconcile everything.

## Layout

- **main.go** — the `Notifier` interface, its actor-backed implementation and
  `fx` wiring.
- **types.go** — the event model and the messages exchanged with the actor.
- **worker.go** — the actor body that owns the subscriber registry.
- **mocks/** — a testify mock of `Notifier` for tests.

## Delivery model

`Subscribe` blocks until the actor confirms the registration, so no event can
slip through between subscribing and the first delivery. `Notify`, by contrast,
returns as soon as the event is queued: delivery to subscribers happens
asynchronously and best-effort.
