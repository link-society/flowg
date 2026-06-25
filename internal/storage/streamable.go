package storage

import (
	"context"
	"io"
)

// Streamable is implemented by stores whose full contents can be dumped to and
// loaded from a byte stream, enabling backup and restore.
type Streamable interface {
	// Dump writes a backup of the store to w, including only changes newer than
	// the given version, and returns the version reached.
	Dump(context.Context, io.Writer, uint64) (uint64, error)
	// Load restores the store from a backup previously produced by Dump.
	Load(context.Context, io.Reader) error
}
