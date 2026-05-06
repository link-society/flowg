---
sidebar_position: 3
---

# What Is a Transformer?

A Transformer is a [VRL script](https://vector.dev/docs/reference/vrl/) used to
transform a log record into another.

They are used to parse, refine, and enrich log records within a pipeline.

A common use case is to parse log records comming from **syslog**:

```vrl
. = parse_syslog!(.message)
```

The above script expects a record with a single field `message`, containing text
in the **syslog** format. After processing, the resulting record will have
fields such as `program`, `timestamp`, or `hostname`.

A transformer can emit multiple logs from a single input by returning an array:

```vrl
. = [
  {"message": "first"},
  {"message": "second"}
]
```

The above script would emit 2 log records with the field `message`.

All datatypes support by VRL are also normalized after the script execution:

```vrl
. = [
  1,
  2,
  "hello",
  {"foo": {"bar": "baz}}
]
```

The above script would emit the following logs:

```json
{"value": "1"}
{"value": "2"}
{"value": "hello"}
{"foo.bar": "baz"}
```
