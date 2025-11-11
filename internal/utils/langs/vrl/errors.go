package vrl

type CompileError struct {
	Message string
}

var _ error = (*CompileError)(nil)

func (e CompileError) Error() string {
	return e.Message
}

type EvalError struct {
	Message string
}

var _ error = (*EvalError)(nil)

func (e EvalError) Error() string {
	return e.Message
}
