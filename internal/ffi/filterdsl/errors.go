package filterdsl

type CompilationError struct {
	Message string
}

func (e *CompilationError) Error() string {
	return e.Message
}

type UnmarshalError struct {
	Reason error
}

func (e *UnmarshalError) Error() string {
	return e.Reason.Error()
}
