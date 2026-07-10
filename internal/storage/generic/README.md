# generic

The packages under `internal/storage/generic` provide reusable, backend-agnostic
building blocks that the domain storage implementations are written against.

They exist so the logic in [databases](../databases) can be expressed once, in
terms of a small abstract key-value API, and then run on top of any concrete
backend that implements it. A backend only has to satisfy the generic contracts
here; it never has to re-implement the domain logic.

## Layout

- **[kv](kv)** — the generic key-value store abstraction: composite keys, values,
  read-only and read-write transactions, and the `Adapter` that ties them to a
  concrete database together with snapshot (backup/restore) support.
