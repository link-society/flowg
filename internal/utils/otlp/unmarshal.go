package otlp

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"

	"link-society.com/flowg/internal/models"
)

// UnmarshalJSON decodes an OTLP/HTTP logs export request encoded as protobuf
// JSON and converts it into FlowG log records.
func UnmarshalJSON(body []byte) ([]*models.LogRecord, error) {
	message := collectlogs.ExportLogsServiceRequest{}
	err := protojson.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}

	logRecords := convertToLogRecords(&message)
	return logRecords, nil
}

// UnmarshalProtobuf decodes an OTLP/HTTP logs export request encoded as binary
// protobuf and converts it into FlowG log records.
func UnmarshalProtobuf(body []byte) ([]*models.LogRecord, error) {
	message := collectlogs.ExportLogsServiceRequest{}
	err := proto.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}

	logRecords := convertToLogRecords(&message)
	return logRecords, nil
}
