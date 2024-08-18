package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Fields    map[string]string `json:"fields"`
}

func NewLogEntry(fields map[string]string) *LogEntry {
	return &LogEntry{
		Timestamp: time.Now(),
		Fields:    fields,
	}
}

func NewLogEntryFromRaw(rawData []byte) (*LogEntry, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, &UnmarshalError{Reason: err}
	}

	flattenedData := map[string]string{}
	flatten("", data, flattenedData)

	return NewLogEntry(flattenedData), nil
}

func (e *LogEntry) NewDbKey(stream string) []byte {
	return []byte(fmt.Sprintf(
		"entry:%s:%020d:%s",
		stream,
		e.Timestamp.UnixMilli(),
		uuid.New().String(),
	))
}
