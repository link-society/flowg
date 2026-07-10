# transactions

The package at `internal/storage/databases/config/transactions` holds the
generic read/write operations behind the config storage. Each function takes a
transaction and an `itemType` selecting the namespace to operate on; this
document describes how that key space is laid out.

## Key space

Every configuration object is stored under a uniform `<itemType>:<name>` key
whose value is the already-serialized item. The same CRUD helpers serve all item
types:

- `pipeline:<name>` — a pipeline definition.
- `transformer:<name>` — a transformer (VRL) definition.
- `forwarder:<name>` — a forwarder definition.
- `system:config` — the single system-wide configuration document.
