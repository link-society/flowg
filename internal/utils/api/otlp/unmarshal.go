package otlp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
