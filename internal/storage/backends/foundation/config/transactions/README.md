# transactions

The package at `internal/storage/backends/foundation/concrete/config/transactions`
holds the generic read/write operations behind the config backend. Each function
takes a FoundationDB transaction and an `itemType` selecting the namespace to
operate on; this document describes how that key space is laid out.

## Key space

FlowG's FoundationDB backend uses a global key space for all key/value pairs it
stores. A subspace named `config` is used for the config storage. Then, for
each item type, another subspace within the `config` subspace is used.

We then pack those subspace prefixes with the item name, and use this as a key.

The value is the already-serialized item.
