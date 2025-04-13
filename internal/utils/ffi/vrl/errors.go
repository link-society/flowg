package vrl

type NullPointerError struct{}

var _ error = (*NullPointerError)(nil)

func (e NullPointerError) Error() string {
	return "received null pointer"
}

type RuntimeError struct {
	Message string
}

var _ error = (*RuntimeError)(nil)

func (e RuntimeError) Error() string {
	return e.Message
}
