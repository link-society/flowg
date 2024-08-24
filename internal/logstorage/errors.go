package logstorage

import "fmt"

type MarshalError struct {
	Reason error
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("failed to marshal entry: %v", e.Reason)
}

type UnmarshalError struct {
	Reason error
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("failed to unmarshal entry: %v", e.Reason)
}

type PersistError struct {
	Operation string
	Reason    error
}

func (e *PersistError) Error() string {
	return fmt.Sprintf("failed to perform persist operation '%s': %v", e.Operation, e.Reason)
}

type QueryError struct {
	Operation string
	Reason    error
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("failed to perform query operation '%s': %v", e.Operation, e.Reason)
}
