package vrl

// CompileError is returned by [NewScriptRunner] when a VRL source fails to
// compile.
type CompileError struct {
	Message string
}

var _ error = (*CompileError)(nil)

// Error implements the error interface.
func (e CompileError) Error() string {
	return e.Message
}

// EvalError is returned by [ScriptRunner.TransformLog] when a VRL program fails
// at runtime.
type EvalError struct {
	Message string
}

var _ error = (*EvalError)(nil)

// Error implements the error interface.
func (e EvalError) Error() string {
	return e.Message
}
