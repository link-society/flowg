package otlp

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"

	"link-society.com/flowg/internal/models"
)

type ContentType string

const (
	ProtoContentType ContentType = "application/x-protobuf"
	JsonContentType  ContentType = "application/json"
)

func UnmarshalLogRecords(body []byte, contentType ContentType) ([]*models.LogRecord, error) {
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
				logRecordModel := models.NewFromOTLP(logRecord)
				logRecords = append(logRecords, logRecordModel)
			}
		}
	}

	return logRecords, nil
}
