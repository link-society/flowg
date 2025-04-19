package otlp

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"

	"link-society.com/flowg/internal/models"
)

func UnmarshalLogRecords(body []byte, contentType string) ([]*models.LogRecord, error) {
	message := collectlogs.ExportLogsServiceRequest{}

	switch contentType {
	case "application/x-protobuf", "application/protobuf":
		err := proto.Unmarshal(body, &message)
		if err != nil {
			return nil, err
		}

	case "application/json":
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
