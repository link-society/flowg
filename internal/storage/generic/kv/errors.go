package kv

import "errors"

var (
	// ErrNotSupported is returned by adapters for operations a particular
	// backend does not implement.
	ErrNotSupported = errors.New("operation not supported")
)
