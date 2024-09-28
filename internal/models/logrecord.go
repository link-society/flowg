package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type LogRecord struct {
	Timestamp time.Time         `json:"timestamp"`
	Fields    map[string]string `json:"fields"`
}

func NewLogRecord(fields map[string]string) *LogRecord {
	return &LogRecord{
		Timestamp: time.Now(),
		Fields:    fields,
	}
}

func (e *LogRecord) NewDbKey(stream string) []byte {
	return []byte(fmt.Sprintf(
		"entry:%s:%020d:%s",
		stream,
		e.Timestamp.UnixMilli(),
		uuid.New().String(),
	))
}
