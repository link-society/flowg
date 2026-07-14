package log

import (
	"fmt"

	"github.com/google/uuid"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"
)

// newDbKey builds the time-ordered storage key for the record in a stream:
// "entry:<stream>:<unix-millis, 20-digit zero-padded>:<uuid>". The padding makes
// a lexical scan walk the stream in chronological order and the uuid keeps
// same-millisecond records distinct.
func newDbKey(stream string, logRecord *models.LogRecord) kv.Key {
	return kv.Key{
		"entry",
		stream,
		fmt.Sprintf("%020d", logRecord.Timestamp.UnixMilli()),
		uuid.New().String(),
	}
}
