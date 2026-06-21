package storage

import (
	"context"
	"io"

	"link-society.com/flowg/internal/storage/changefeed"
)

type Streamable interface {
	Dump(context.Context, io.Writer, uint64) (uint64, error)
	Load(context.Context, io.Reader) error
	Merge(context.Context, io.Reader) error
	ApplyReplicated(context.Context, []changefeed.Record) error
}
