package models

import (
	"strconv"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	"github.com/swaggest/jsonschema-go"
)

type DynamicField string

func (DynamicField) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithPattern(`^@expr:`)
	return nil
}

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
