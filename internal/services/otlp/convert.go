package otlp

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	logsmodels "go.opentelemetry.io/proto/otlp/logs/v1"
	metricsmodels "go.opentelemetry.io/proto/otlp/metrics/v1"
	tracesmodels "go.opentelemetry.io/proto/otlp/trace/v1"
	"link-society.com/flowg/internal/models"
)

// ToFields convert a struct to a map[string]string
func ToFields(obj interface{}) (map[string]string, error) {
	out := make(map[string]string)
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %T", obj)
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name
		tag := fieldType.Tag.Get("json")
		if tag != "" && tag != "-" {
			if commaIndex := len(tag); commaIndex > 0 && tag[commaIndex-1] == ',' {
				tag = tag[:commaIndex-1]
			}
			fieldName = tag
		}

		b, err := json.Marshal(field.Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field %s: %w", fieldName, err)
		}
		out[fieldName] = string(b)
	}
	return out, nil
}

func MetricToLogRecord(metric *metricsmodels.Metric) (result models.LogRecord, err error) {
	result = models.LogRecord{
		Timestamp: time.Now(),
	}
	result.Fields, err = ToFields(metric)

	return
}

func LogToLogRecord(logRecord *logsmodels.LogRecord) (result models.LogRecord, err error) {
	result = models.LogRecord{
		Timestamp: time.Unix(0, int64(logRecord.TimeUnixNano)),
	}

	result.Fields, err = ToFields(logRecord)

	return
}

func SpanToLogRecord(span *tracesmodels.Span) (result models.LogRecord, err error) {
	result = models.LogRecord{
		Timestamp: time.Unix(0, int64(span.StartTimeUnixNano)),
	}

	result.Fields, err = ToFields(span)

	return
}
