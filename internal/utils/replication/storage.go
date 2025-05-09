package replication

import (
	"context"
	"io"
)

type Storage interface {
	Dump(context.Context, io.Writer, uint64) error
	Load(context.Context, io.Reader) error
}
