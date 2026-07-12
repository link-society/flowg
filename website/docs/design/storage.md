---
sidebar_position: 1
---

# How Logs Are Stored?

**FlowG** stores everything in a **composite key/value** model: each key is an
ordered list of string segments, and each value is an opaque blob.

This logical model is backend-agnostic. **FlowG** ships with two storage
backends that implement it:

 - [BadgerDB](https://github.com/dgraph-io/badger) (the default): an embedded,
   single-node key/value store.
 - [FoundationDB](https://www.foundationdb.org/): a distributed, horizontally
   scalable key/value store.

Throughout the sections below, keys are written using the **BadgerDB**
representation, where the segments of a composite key are joined with a colon
(`:`). For example, the composite key `["entry", "test", ...]` is written
`entry:test:...`. See [FoundationDB Backend](#foundationdb-backend) for how the
same keys are mapped when using FoundationDB.

## Streams

A stream is the destination of log entries. The stream hold the configuration
for the following:

 - **Time-based retention:** how long will the log entry be kept (in seconds)
 - **Size-based retention:** the maximum size (in MB) of the stream before deleting older entries
 - **Indexed fields:** the list of fields in the log entries to index

If any of the retention parameters is set to 0 (the default), it is assumed to
be "unlimited". For each stream, the configuration is stored (as JSON) at the
following key:

```
stream:config:<stream name> = { ... }
```

Example:

```
# TTL of 5min and max size of 1GB
stream:config:test = {"ttl": 300, "size": 1024, "indexed_fields": ["foo", "bar"]}
```

## Log Entries

A log entry is composed of:

 - the timestamp at which the log was ingested
 - a flat record of string key/values

Each log entry is serialized as JSON and saved at the following key:

```
entry:<stream name>:<timestamp in milliseconds>:<UUID v4>
```

The timestamp is padded with 0s to be easily sortable. The UUID is there to make
sure that 2 entries can be ingested at the same millisecond.

Example:

```
entry:test:00000001724140167373:057804d1-832f-45bf-8e70-7acbf22ec480
```

When added to the database, if no stream configuration exists, a default one is
added.

For each field of the log record, the following key is added to the database:

```
stream:field:<stream name>:<field name>
```

Example:

```
stream:field:test:foo
stream:field:test:bar
```

## Indexes

### Field Index

A field of a log entry is indexed by creating a new key referencing the field
and the associated log entry's key:

```
index:<stream name>:field:<field name>:<base64 encoded field value>:<entry key>
```

Example:

```
# base64(bar) = YmFy
index:test:field:foo:YmFy:entry:test:00000001724140167373:057804d1-832f-45bf-8e70-7acbf22ec480
```

### Index value size limit

The field value lives *inside* the index key (base64-encoded), and every backend
caps the size of a key — FoundationDB rejects any key larger than 10&nbsp;KB. If
an indexed value is large enough to push the index key past that limit, the
record is still stored, but that value is **not** indexed (the write is skipped
with a logged warning rather than failing ingestion). Such records are returned
by time-range queries but never match an exact-value filter on that field.

The threshold is deterministic. The index key is

```
index / <stream> / field / <field> / <base64(value)> / entry / <stream> / <20-digit millis> / <36-char uuid>
```

so its encoded size grows with the stream name (which appears twice), the field
name, and the base64-encoded value (≈ 4⁄3 of the raw value). Taking the stricter
FoundationDB backend (tuple-encoded key inside the `flowg/log` subspace), a raw
field value of `V` bytes is indexable as long as

```
101 + 2·len(stream) + len(field) + 4·⌈V/3⌉ ≤ 10000
```

which gives the maximum indexable value size

```
V_max = 3 · ⌊(9899 − 2·len(stream) − len(field)) / 4⌋   bytes
```

For example, stream `default` (7 chars) and field `level` (5 chars) yield
`V_max ≈ 7410` bytes (~7.2&nbsp;KiB). Backends never rely on this estimate: each
one checks the *actual* encoded key size, so the real threshold is always exact
for the backend in use.

## Querying

Each query must be made within a time-window.

We fetch the keys to lookup by iterating over all keys with the following
prefix:

```
entry:<stream name>:
```

Then, we select all keys that are "higher" than

```
entry:<stream name>:<from timestamp>:
```

And "lower" than:

```
entry:<stream name>:<to timestamp>:
```

Because of the internal structure of *BadgerDB*, this operation is fast.

Then, if a [filter](/docs/user/guides/filtering) is given, we match the log
record to the filter expression to determine if it should be returned.

## FoundationDB Backend

When configured to use **FoundationDB**, the exact same composite keys described
above are mapped onto FoundationDB's native key space:

 - Each composite key is packed into a
   [tuple](https://apple.github.io/foundationdb/data-modeling.html#tuples)
   inside a
   [subspace](https://apple.github.io/foundationdb/developer-guide.html#subspaces)
   scoped to `<keyspace>/<namespace>`. Logs live in the `flowg/log` subspace
   (`flowg` being the default keyspace).
 - The subspace prefix is transparently prepended on write and stripped on read,
   so it never leaks into the logical keys. In other words, the code still works
   with `entry:<stream>:...`, `index:<stream>:...`, etc.
 - Tuple packing preserves the lexicographic ordering of the segments, so the
   time-window range scans used by [Querying](#querying) remain fast.

### Time-based retention

Unlike BadgerDB, FoundationDB has no native TTL. Instead, every value is wrapped
in an envelope whose first 8 bytes hold the expiration timestamp (unix seconds,
big-endian; `0` meaning "no expiration"):

```
<8-byte expiry timestamp><payload>
```

Expiration is then enforced in two complementary ways:

 - **Lazily on read:** expired entries are skipped when reading or iterating, so
   they are never returned even before being physically removed.
 - **Eagerly via garbage collection:** a background collector periodically scans
   for expired keys and deletes them to reclaim space.
