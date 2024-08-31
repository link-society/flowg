# How Logs Are Stored?

The storage backend of **FlowG** is the key/value store
[BadgerDB](https://dgraph.io/docs/badger/).

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

## Indexes

### Stream Index

For each stream, an empty value is stored at the following key:

```
stream:<stream name>
```

Example:

```
stream:test
```

### Field Index

Each field of a log entry is indexed by creating a new key referencing the
field and the associated log entry's key:

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

Then, if a [filter](../guides/filtering.md) is given:

For `foo = "bar"` filters, we fetch all keys within the time-window with the
following prefix:

```
index:<stream name>:field:foo:YmFy:
```

For `foo in ["bar", "baz"]` filters, we fetch all keys the time-window with
the following prefixes:

```
# base64(bar) = YmFy
# base64(baz) = YmF6

index:<stream name>:field:foo:YmFy:
index:<stream name>:field:foo:YmF6:
```

For `not <sub-filter>` filters, we select all keys in the time-window that do
not appear in the `<sub-filter>` result.

> **NB:**
>  - `field != value` is desugared into `not (field = value)`
>  - `field not in [value, value]` is desugared into `not (field in [value, value])`

For `<lhs> or <rhs>` filters, we select the union of `<lhs>` and `<rhs>`
results.

For `<lhs> and <rhs>` filters, we select the intersection of `<lhs>` and `<rhs>`
results.
