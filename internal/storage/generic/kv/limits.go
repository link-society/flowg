package kv

import "fmt"

// FoundationDB enforces hard limits on the byte size of the keys and values a
// transaction may write, and exceeding them aborts the whole commit with an
// opaque error. FlowG adopts these limits as part of the kv contract so every
// backend agrees on what it will accept and rejects oversized mutations up front
// with a typed error ([ErrKeyTooLarge] / [ErrValueTooLarge]) instead of failing
// deep inside an engine.
//
// See https://apple.github.io/foundationdb/known-limitations.html.
const (
	// MaxKeySize is the maximum size, in bytes, of an encoded key.
	MaxKeySize = 10_000
	// MaxValueSize is the maximum size, in bytes, of an encoded value.
	MaxValueSize = 100_000
)

// CheckKeySize reports whether an encoded key of size bytes fits within
// [MaxKeySize], returning a wrapped [ErrKeyTooLarge] otherwise.
//
// Backends pass the size of the key as it will be stored (including any prefix
// or envelope they add), so the check reflects what the engine actually sees.
func CheckKeySize(size int) error {
	if size > MaxKeySize {
		return fmt.Errorf("%w: %d bytes exceeds limit of %d", ErrKeyTooLarge, size, MaxKeySize)
	}
	return nil
}

// CheckValueSize reports whether an encoded value of size bytes fits within
// [MaxValueSize], returning a wrapped [ErrValueTooLarge] otherwise.
//
// Backends pass the size of the value as it will be stored (including any
// envelope they add), so the check reflects what the engine actually sees.
func CheckValueSize(size int) error {
	if size > MaxValueSize {
		return fmt.Errorf("%w: %d bytes exceeds limit of %d", ErrValueTooLarge, size, MaxValueSize)
	}
	return nil
}
