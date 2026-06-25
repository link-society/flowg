# lognotify

The package at `internal/engines/lognotify` is the live notification bus for
ingested log records. It lets clients subscribe to a stream and receive every
record routed to it from that moment on, which is what powers the live tail in
the web UI.

It is built on a single [actor](https://github.com/vladopajic/go-actor): all
subscriber bookkeeping happens on one goroutine, so the per-stream subscriber
sets need no locking.

## Responsibilities

- **Subscriptions** — registers per-subscriber mailboxes against a stream and
  tears them down automatically when the subscriber's context is cancelled.
- **Fan-out** — broadcasts each notified record to every current subscriber of
  its stream, logging — but never failing on — individual delivery errors.
- **Wiring** — `NewLogNotifier` returns an `fx` module providing a `LogNotifier`
  whose actor and mailboxes follow the application lifecycle.

## Layout

- **main.go** — the `LogNotifier` interface, its actor-backed implementation and
  `fx` wiring.
- **types.go** — the messages exchanged with the actor (`LogMessage`,
  `SubscribeMessage`, `ReadyResponse`).
- **worker.go** — the actor body that owns the subscriber registry.
- **mock.go** — a testify mock of `LogNotifier` for tests.

## Delivery model

`Subscribe` blocks until the actor confirms the registration, so no record can
slip through between subscribing and the first delivery. `Notify`, by contrast,
returns as soon as the record is queued: delivery to subscribers happens
asynchronously and best-effort.
