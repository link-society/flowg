package filterdsl

type CompilationError struct {
	Message string
}

var _ error = (*CompilationError)(nil)

func (e *CompilationError) Error() string {
	return e.Message
}

type UnmarshalError struct {
	Reason error
}

var _ error = (*UnmarshalError)(nil)

func (e *UnmarshalError) Error() string {
	return e.Reason.Error()
}
