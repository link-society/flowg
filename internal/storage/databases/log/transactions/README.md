# transactions

The package at `internal/storage/databases/log/transactions` holds the
low-level read/write operations behind the log storage. Each function takes a
transaction and manipulates the log key space directly; this document describes
how that key space is laid out.

## Key space

A stream's records, its configuration, its field set and its inverted indices
each live under their own prefix, so a stream is the union of the keys sharing
its name.

### Records

- `entry:<stream>:<timestamp>:<id>` — one key per ingested record; its value is
  the JSON-encoded `LogRecord`. The timestamp is zero-padded so keys sort
  chronologically, which lets time-range queries seek to a bound and stop early.
  Records carry the stream's retention TTL and expire on their own.

### Stream configuration and fields

- `stream:config:<stream>` — the JSON-encoded stream configuration; an empty
  value is treated as the default (zero-value) config, and its mere existence
  marks the stream.
- `stream:field:<stream>:<field>` — one existence marker per field name ever seen
  in the stream, used to enumerate a stream's fields.

### Inverted indices

- `index:<stream>:field:<field>:<base64(value)>:<entry-key...>` — one existence
  marker per (field, value) pair pointing back at a record, created only for the
  fields a stream is configured to index. The value segment is base64-encoded so
  it cannot collide with the `:` key separator, and index keys inherit the
  referenced entry's TTL so the two expire together. Because the value lives
  inside the key, a value large enough to push the index key past the backend's
  key-size limit (`kv.MaxKeySize`) cannot be indexed: `AddKey` catches
  `kv.ErrKeyTooLarge`, logs a warning and skips that value, so ingestion still
  succeeds and the record stays queryable by time — it just won't match an
  exact-value filter on that field.

## Notes

- Ingesting a record writes its `entry:` key, registers each of its fields, and,
  for every indexed field, adds the matching `index:` key — all in one
  transaction.
- Indexed-field queries intersect the time-window candidates with the `index:`
  keys before decoding any record; `Distinct` reads values straight from the
  index keys without touching records at all.
- Retention is enforced two ways: per-record TTLs (retention time) and a garbage
  collector that evicts the oldest records once a stream exceeds its retention
  size, purging their index references as it goes.
