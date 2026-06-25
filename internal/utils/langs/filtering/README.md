# filtering

The package at `internal/utils/langs/filtering` provides the boolean expression
language used to filter log records.

It exists to let queries select records with arbitrary predicates over their
fields while keeping the evaluation logic in one place. A query string is
compiled once into a reusable [Filter] and then evaluated against many records,
so the parsing cost is paid only once per query.

## Responsibilities

- **Compilation** — `Compile` parses a filter expression and compiles it into a
  `Filter`, returning an error for malformed expressions.
- **Evaluation** — `Filter.Evaluate` reports whether a given log record matches
  the compiled expression.

## Language

Expressions are written in the [expr](https://expr-lang.org/) language and must
evaluate to a boolean. A record's fields are exposed to the expression as
variables; referencing an undefined field is allowed and yields a nil value
rather than an error, so filters degrade gracefully across heterogeneous
records.
