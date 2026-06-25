# transactions

The package at `internal/storage/backends/badger/concrete/log/transactions`
holds the low-level read/write, indexing and retention operations behind the log
backend. Each function takes a Badger transaction and manipulates the log key
space directly; this document describes how that key space is laid out.

## Key space

### Records

- `entry:<stream>:<unix-millis>:<uuid>` — one JSON-encoded log record. The
  timestamp is zero-padded to 20 digits, so a lexical scan of an
  `entry:<stream>:` prefix walks the stream in chronological order, and a query
  can seek straight to a time range instead of scanning the whole stream.

### Stream metadata

- `stream:config:<stream>` — the JSON-encoded stream configuration (retention,
  indexed fields). An empty value is treated as "defaults".
- `stream:field:<stream>:<field>` — an existence marker for every field name ever
  seen in the stream, used to advertise its queryable fields.

### Inverted index

- `index:<stream>:field:<field>:<base64(value)>:<entryKey>` — maps a field value
  back to the records that carry it. The value is base64-encoded so arbitrary
  field contents can live safely inside the colon-delimited key, and the trailing
  entry key points at the matching record.

## Notes

- Index keys are given the same TTL as the entries they reference, so they expire
  together.
- A query combines up to three filters: a time range (expressed through the key
  range), an optional field-index lookup (values OR-ed within a field, fields
  AND-ed together), and an optional [expr](https://expr-lang.org/) filter
  evaluated against the decoded record.
- Retention is enforced two ways: time-based retention relies on Badger TTLs set
  on entries at ingestion, while size-based retention is applied explicitly by
  the garbage collector, which evicts the oldest records of any over-budget
  stream until it fits within its configured size.
- In demo mode the stream configuration is never read from disk: a fixed
  in-memory config is returned for every stream.
