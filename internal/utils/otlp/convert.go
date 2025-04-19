package otlp

import (
	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"

	"link-society.com/flowg/internal/models"
)

func convertToLogRecords(message *collectlogs.ExportLogsServiceRequest) []*models.LogRecord {
	var logRecords []*models.LogRecord

	for _, resourceLogs := range message.GetResourceLogs() {
		for _, scopeLogs := range resourceLogs.GetScopeLogs() {
			for _, logRecord := range scopeLogs.GetLogRecords() {
				logRecordModel := models.NewFromOTLP(logRecord)
				logRecords = append(logRecords, logRecordModel)
			}
		}
	}

	return logRecords
}
