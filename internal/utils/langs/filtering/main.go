package filtering

// Compile parses and compiles a filter expression into a reusable [Filter].
//
// The expression is written in the [expr] language and must evaluate to a
// boolean. It is compiled once here and can then be evaluated against many log
// records, so callers should compile a query a single time and reuse the
// returned [Filter].
//
// [expr]: https://expr-lang.org/
func Compile(input string) (Filter, error) {
	return newFilterImpl(input)
}
