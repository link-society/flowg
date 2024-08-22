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

Each field of a log entry is indexed by adding the log entry's key in a list,
stored at the following key:

```
index:<stream name>:field:<field name>:<field value>
```

Example:

```
index:test:field:foo:bar = [
  entry:test:00000001724140167373:057804d1-832f-45bf-8e70-7acbf22ec480
]
```

## Querying

Each query must be made within a time-window.

We fetch the keys to lookup by iterating over all keys with the following
prefix:

```
stream:<stream name>:
```

Then, we select all keys that are "higher" than

```
stream:<stream name>:<from timestamp>:
```

And "lower" than:

```
stream:<stream name>:<to timestamp>:
```

Because of the internal structure of *BadgerDB*, this operation is fast.

Then, if a [filter](../guides/filtering.md) is given, we fetch the indexes that
match the filter by iterating over keys with the following prefixes:

```
index:<stream name>:field:<field name>:
```

Once we got the complete list of log entry keys to fetch, we sort them and fetch
the actual records.
