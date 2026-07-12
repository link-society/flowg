---
sidebar_position: 2
---

# What Are Streams?

A stream is the destination of a log record after being processed by a pipeline.
This is the actual collection that you will query/visualize using
[filters](/docs/user/guides/filtering).

Every field can be indexed, allowing fast retrieval of log records matching some
[filter](/docs/user/guides/filtering).

:::warning Large field values are not indexed

A field's value is embedded in its index key, and storage backends cap the size
of a key (FoundationDB rejects keys larger than 10&nbsp;KB). When an indexed
field's value is too large to fit in an index key, FlowG **stores the record but
skips indexing that value** — it does not reject the record.

Such a record is still returned by time-range queries, but it will **not** match
an exact-value [filter](/docs/user/guides/filtering) on that field, and its value
will not appear in the field's value suggestions.

With typical (short) stream and field names, values up to roughly **7&nbsp;KB**
are always indexed. The exact threshold is deterministic but depends on the
backend and on the lengths of the stream and field names — see the
[storage design](/docs/design/storage) for the formula. Avoid indexing
free-form, high-cardinality fields such as message bodies.

:::

