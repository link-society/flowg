package hash

type HashError struct {
	Reason error
}

func (e *HashError) Error() string {
	return e.Reason.Error()
}
