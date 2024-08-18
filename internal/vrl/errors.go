package vrl

type NullPointerError struct{}

func (e NullPointerError) Error() string {
	return "received null pointer"
}

type RuntimeError struct {
	Message string
}

func (e RuntimeError) Error() string {
	return e.Message
}
