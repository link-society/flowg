package otlp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	otlpmodels "go.opentelemetry.io/proto/otlp/logs/v1"
	"google.golang.org/protobuf/proto"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"

	"link-society.com/flowg/internal/models"
)

type ContentType string

const (
	ProtoContentType ContentType = "application/x-protobuf"
	JsonContentType  ContentType = "application/json"
)

func UnmarshalLogRecords(r *http.Request) ([]*models.LogRecord, error) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	contentType := ContentType(r.Header.Get("Content-Type"))
	message := collectlogs.ExportLogsServiceRequest{}

	switch contentType {
	case ProtoContentType:
		err := proto.Unmarshal(body, &message)
		if err != nil {
			return nil, err
		}

	case JsonContentType:
		err := json.Unmarshal(body, &message)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	var logRecords []*models.LogRecord

	for _, resourceLogs := range message.GetResourceLogs() {
		for _, scopeLogs := range resourceLogs.GetScopeLogs() {
			for _, logRecord := range scopeLogs.GetLogRecords() {
				logRecordModel, err := LogToLogRecord(logRecord)
				if err != nil {
					return nil, err
				}
				logRecords = append(logRecords, logRecordModel)
			}
		}
	}

	return logRecords, nil
}

func OTLPDataToLogRecordFields(obj interface{}) (map[string]string, error) {
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

func LogToLogRecord(logRecord *otlpmodels.LogRecord) (result *models.LogRecord, err error) {
	result = &models.LogRecord{
		Timestamp: time.Unix(0, int64(logRecord.TimeUnixNano)),
	}

	result.Fields, err = OTLPDataToLogRecordFields(logRecord)

	return
}
