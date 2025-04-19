package otlp

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"

	"link-society.com/flowg/internal/models"
)

func UnmarshalJSON(body []byte) ([]*models.LogRecord, error) {
	message := collectlogs.ExportLogsServiceRequest{}
	err := json.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}

	logRecords := convertToLogRecords(&message)
	return logRecords, nil
}

func UnmarshalProtobuf(body []byte) ([]*models.LogRecord, error) {
	message := collectlogs.ExportLogsServiceRequest{}
	err := proto.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}

	logRecords := convertToLogRecords(&message)
	return logRecords, nil
}
