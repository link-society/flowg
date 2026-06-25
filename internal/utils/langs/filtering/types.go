package filtering

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	"link-society.com/flowg/internal/models"
)

// Filter is a compiled boolean predicate over a log record, produced by
// [Compile]. It decides whether a record matches a query.
type Filter interface {
	// Evaluate reports whether record satisfies the filter expression. The
	// record's fields are exposed to the expression as variables; referencing an
	// undefined field is allowed and yields a nil value rather than an error.
	Evaluate(record *models.LogRecord) (bool, error)
}

type filterImpl struct {
	program *vm.Program
}

var _ Filter = (*filterImpl)(nil)

func newFilterImpl(input string) (Filter, error) {
	program, err := expr.Compile(
		input,
		expr.Env(map[string]string{}),
		expr.AllowUndefinedVariables(),
		expr.AsBool(),
		expr.WarnOnAny(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to compile expression: %w", err)
	}

	return &filterImpl{program: program}, nil
}

func (f *filterImpl) Evaluate(record *models.LogRecord) (bool, error) {
	output, err := expr.Run(f.program, record.Fields)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate expression: %w", err)
	}

	return output.(bool), nil
}
