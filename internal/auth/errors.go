package auth

import "fmt"

type InvalidScopeError struct {
	Scope string
}

func (e *InvalidScopeError) Error() string {
	return fmt.Sprintf("invalid scope: %s", e.Scope)
}

type PersistError struct {
	Operation string
	Reason    error
}

func (e *PersistError) Error() string {
	return fmt.Sprintf("failed to perform persist operation '%s': %s", e.Operation, e.Reason)
}

type QueryError struct {
	Operation string
	Reason    error
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("failed to perform query operation '%s': %s", e.Operation, e.Reason)
}
