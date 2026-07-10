package kv

import (
	"iter"
	"time"
)

// Represents operations on a read-only transaction
type QueryTx interface {
	// Retrieve the value associated with the given key.
	Get(key Key) (Value, error)
	// Iterate over keys with the given prefix.
	IterKeys(prefix Key, keyRange KeyRange) iter.Seq[Key]
	// Iterate over key-value pairs with the given prefix.
	IterPairs(prefix Key, keyRange KeyRange) iter.Seq[Pair]
}

// Represents operations on a read-write transaction
type MutationTx interface {
	QueryTx

	// Set the value for the given key.
	Set(key Key, value Value) error
	// Set the value for the given key with a TTL (time-to-live) duration.
	SetWithTTL(key Key, value Value, ttl time.Duration) error
	// Delete the key-value pair associated with the given key.
	Clear(key Key) error
}
