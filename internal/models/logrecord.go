package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	otlpcommonmodels "go.opentelemetry.io/proto/otlp/common/v1"
	otlplogmodels "go.opentelemetry.io/proto/otlp/logs/v1"
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

// NewFromOTLP flattens an OpenTelemetry log record into a LogRecord, promoting
// the well-known OTLP fields to top-level names and prefixing each attribute
// with "attr.".
func NewFromOTLP(logRecord *otlplogmodels.LogRecord) *LogRecord {
	fields := map[string]string{
		"severity_number":          logRecord.SeverityNumber.String(),
		"severity_text":            logRecord.SeverityText,
		"body":                     otlpValueParser(logRecord.Body),
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
		fields[fieldName] = otlpValueParser(attribute.Value)
	}

	return &LogRecord{
		Timestamp: time.Now(),
		Fields:    fields,
	}
}

// NewDbKey builds the time-ordered storage key for the record in a stream:
// "entry:<stream>:<unix-millis, 20-digit zero-padded>:<uuid>". The padding makes
// a lexical scan walk the stream in chronological order and the uuid keeps
// same-millisecond records distinct.
func (e *LogRecord) NewDbKey(stream string) []byte {
	return []byte(fmt.Sprintf(
		"entry:%s:%020d:%s",
		stream,
		e.Timestamp.UnixMilli(),
		uuid.New().String(),
	))
}

// otlpValueParser renders an arbitrary OTLP AnyValue as a string so it can live
// in a LogRecord's flat field map.
func otlpValueParser(v *otlpcommonmodels.AnyValue) string {
	switch v.Value.(type) {
	case *otlpcommonmodels.AnyValue_StringValue:
		return v.GetStringValue()
	case *otlpcommonmodels.AnyValue_BoolValue:
		return fmt.Sprintf("%t", v.GetBoolValue())
	case *otlpcommonmodels.AnyValue_IntValue:
		return fmt.Sprintf("%d", v.GetIntValue())
	case *otlpcommonmodels.AnyValue_DoubleValue:
		return fmt.Sprintf("%f", v.GetDoubleValue())
	case *otlpcommonmodels.AnyValue_ArrayValue:
		items := make([]string, 0, len(v.GetArrayValue().Values))

		for i, item := range v.GetArrayValue().Values {
			items[i] = otlpValueParser(item)
		}

		return fmt.Sprintf("[%s]", strings.Join(items, " "))

	case *otlpcommonmodels.AnyValue_KvlistValue:
		items := make([]string, 0, len(v.GetKvlistValue().Values))

		for i, item := range v.GetKvlistValue().Values {
			items[i] = fmt.Sprintf("%s:%s", item.Key, otlpValueParser(item.Value))
		}

		return fmt.Sprintf("map[%s]", strings.Join(items, " "))

	case *otlpcommonmodels.AnyValue_BytesValue:
		return fmt.Sprintf("%v", v.GetBytesValue())

	default:
		return v.String()
	}
}
