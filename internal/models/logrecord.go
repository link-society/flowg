package models

import (
	"time"
)

// LogRecord is the canonical in-memory representation of a log entry throughout
// FlowG: a timestamp plus a flat map of string fields. Every ingestion source
// (text, structured, OTLP, syslog) is normalised into this shape before it
// enters a pipeline.
type LogRecord struct {
	Timestamp time.Time         `json:"timestamp" required:"true" format:"date-time"`
	Fields    map[string]string `json:"fields" required:"true"`
}

// NewLogRecord builds a record from a field map, stamping it with the current
// time.
func NewLogRecord(fields map[string]string) *LogRecord {
	return &LogRecord{
		Timestamp: time.Now(),
		Fields:    fields,
	}
}
