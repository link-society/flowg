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
