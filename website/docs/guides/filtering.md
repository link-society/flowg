---
sidebar_position: 1
---

# Filtering Data

In **FlowG**, every log record is a flat record where every field is a string.
Every field in a [stream](/docs/guides/streams) is indexed, allowing fast
retrieval of log records matching some filter.

Said filter are written using a very small DSL.

## Syntax

| Example | Description |
| --- | --- |
| `field = "value"` | Matches all records where `field` is strictly equal to `value` |
| `field != "value"` | Matches all records where `field` is **not** equal to `value` |
| `field in ["a", "b"]` | Matches all records where `field` is either `a` or `b` |
| `field not in ["a", "b"]` | Matches all records where `field` is neither `a` or `b` |

Filters can be combined in more complex filters using the boolean operators
`not`, `and`, `or`:

| Example | Description |
| --- | --- |
| `field1 = "a" or field2 = "b"` | Matches all records where either `field1` equals `a` or `field2` equals `b` |
| `field1 = "a" and field2 = "b"` | Matches all records where both `field1` equals `a` and `field2` equals `b` |
| `not (field1 = "a" or field2 = "b")` | Matches all records where neither `field1` equals `a` and `field2` equals `b` |
