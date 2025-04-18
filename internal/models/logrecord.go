package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	otlpmodels "go.opentelemetry.io/proto/otlp/logs/v1"
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

func NewFromOTLP(logRecord *otlpmodels.LogRecord) *LogRecord {
	fields := map[string]string{
		"severity_number":          logRecord.SeverityNumber.String(),
		"severity_text":            logRecord.SeverityText,
		"body":                     logRecord.Body.String(),
		"dropped_attributes_count": fmt.Sprintf("%d", logRecord.DroppedAttributesCount),
		"flags":                    fmt.Sprintf("%d", logRecord.Flags),
		"trace_id":                 fmt.Sprintf("%x", logRecord.TraceId),
		"span_id":                  fmt.Sprintf("%x", logRecord.SpanId),
		"event_name":               logRecord.EventName,
		"observed_time_unix_nano":  fmt.Sprintf("%d", logRecord.ObservedTimeUnixNano),
		"time_unix_nano":           fmt.Sprintf("%d", logRecord.TimeUnixNano),
	}
	for _, attribute := range logRecord.Attributes {
		fieldName := fmt.Sprintf("attr.%s", attribute.Key)
		fieldValue := attribute.Value.String()
		fields[fieldName] = fieldValue
	}

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
