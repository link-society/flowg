package kv

import "sort"

// Represents a composite key in a Key-Value database.
type Key []string

// Represents a range of keys for iteration
type KeyRange struct {
	// The start of the range (inclusive), ignored if nil
	From Key
	// The end of the range (inclusive), ignored if nil
	To Key
}

// Represents a sortable sequence of keys
type KeySlice []Key

var _ sort.Interface = KeySlice{}

// Represents arbitrary values in a Key-Value database.
type Value []byte

// Represents a key-value pair in a Key-Value database.
type Pair interface {
	// Fetch the key (with copy)
	Key() Key

	// Fetch the value (with copy)
	Value() Value

	// Return an estimation of the value's size
	EstimateSize() int64

	// Return the expiration time of the value, in seconds since the Unix epoch.
	// If the value does not have an expiration time, return 0.
	ExpiresAt() uint64
}

// HasSuffix reports whether the key ends with the given sequence of segments.
func (k Key) HasSuffix(suffix Key) bool {
	if len(k) < len(suffix) {
		return false
	}

	for i := 0; i < len(suffix); i++ {
		if k[len(k)-len(suffix)+i] != suffix[i] {
			return false
		}
	}

	return true
}

// Len returns the number of keys in the slice, implementing [sort.Interface].
func (ks KeySlice) Len() int {
	return len(ks)
}

// Less reports whether the key at index i sorts before the one at index j,
// comparing segment by segment, implementing [sort.Interface].
func (ks KeySlice) Less(i, j int) bool {
	for x := 0; x < len(ks[i]) && x < len(ks[j]); x++ {
		if ks[i][x] < ks[j][x] {
			return true
		} else if ks[i][x] > ks[j][x] {
			return false
		}
	}

	return len(ks[i]) < len(ks[j])
}

// Swap exchanges the keys at indices i and j, implementing [sort.Interface].
func (ks KeySlice) Swap(i, j int) {
	ks[i], ks[j] = ks[j], ks[i]
}
