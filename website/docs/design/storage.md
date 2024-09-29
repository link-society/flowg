# How Logs Are Stored?

The storage backend of **FlowG** is the key/value store
[BadgerDB](https://dgraph.io/docs/badger/).

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

Then, if a [filter](/docs/guides/filtering) is given:

For `foo = "bar"` filters, we fetch all keys within the time-window with the
following prefix:

```
index:<stream name>:field:foo:YmFy:
```

If the field is not indexed, all keys in the time-window are returned.

For `foo in ["bar", "baz"]` filters, we fetch all keys the time-window with
the following prefixes:

```
# base64(bar) = YmFy
# base64(baz) = YmF6

index:<stream name>:field:foo:YmFy:
index:<stream name>:field:foo:YmF6:
```

If the field is not indexed, all keys in the time-window are returned.

For `not <sub-filter>` filters, we select all keys in the time-window that do
not appear in the `<sub-filter>` result.

> **NB:**
>  - `field != value` is desugared into `not (field = value)`
>  - `field not in [value, value]` is desugared into `not (field in [value, value])`

For `<lhs> or <rhs>` filters, we select the union of `<lhs>` and `<rhs>`
results.

For `<lhs> and <rhs>` filters, we select the intersection of `<lhs>` and `<rhs>`
results.

Once the set of keys to fetch is determined, we load the values and run them
through the filter once again (to account for unindexed fields), and return the
final set of log entries.
