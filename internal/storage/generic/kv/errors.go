package kv

import "errors"

var (
	// ErrNotSupported is returned by adapters for operations a particular
	// backend does not implement.
	ErrNotSupported = errors.New("operation not supported")

	// ErrKeyTooLarge is returned by a mutation when an encoded key exceeds
	// [MaxKeySize].
	ErrKeyTooLarge = errors.New("key too large")

	// ErrValueTooLarge is returned by a mutation when an encoded value exceeds
	// [MaxValueSize].
	ErrValueTooLarge = errors.New("value too large")
)
