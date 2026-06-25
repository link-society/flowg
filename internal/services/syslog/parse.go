package syslog

import (
	"fmt"

	gosyslogformat "gopkg.in/mcuadros/go-syslog.v2/format"

	"link-society.com/flowg/internal/models"
)

// parseLogParts converts the loosely-typed fields parsed by the syslog library
// into a LogRecord, stringifying every value so it fits the flat field map.
func parseLogParts(logParts gosyslogformat.LogParts) *models.LogRecord {
	fields := make(map[string]string, len(logParts))

	for key, value := range logParts {
		switch value := value.(type) {
		case string:
			fields[key] = value

		case []byte:
			fields[key] = string(value)

		default:
			fields[key] = fmt.Sprintf("%v", value)
		}
	}

	return models.NewLogRecord(fields)
}
