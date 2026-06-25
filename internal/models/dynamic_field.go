package models

import (
	"strconv"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	"github.com/swaggest/jsonschema-go"
)

// DynamicField is a forwarder configuration value that may either be a literal
// string or an [expr](https://expr-lang.org/) expression (when prefixed with
// "@expr:"), evaluated against the record at forward time.
type DynamicField string

// PrepareJSONSchema describes a DynamicField for OpenAPI generation: a string
// matching the "@expr:" convention.
func (DynamicField) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.AddType(jsonschema.String)
	schema.WithPattern(`^@expr:`)
	return nil
}

// CompileDynamicField compiles a dynamic field into an expr program. A value
// prefixed with "@expr:" is compiled as an expression; any other value is
// compiled as a quoted string literal, so plain values evaluate to themselves.
func CompileDynamicField(value string) (*vm.Program, error) {
	if len(value) >= 6 && value[:6] == "@expr:" {
		return expr.Compile(
			value[6:],
			expr.Env(map[string]any{}),
			expr.AllowUndefinedVariables(),
		)
	} else {
		return expr.Compile(
			strconv.Quote(value),
			expr.Env(map[string]any{}),
			expr.AllowUndefinedVariables(),
		)
	}
}
