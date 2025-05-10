package storage

import (
	"context"
	"io"
)

type Streamable interface {
	Dump(context.Context, io.Writer, uint64) (uint64, error)
	Load(context.Context, io.Reader) error
}
