package syslog

import (
	"testing"

	"fmt"
	"time"

	gosyslogformat "gopkg.in/mcuadros/go-syslog.v2/format"
)

func TestParseLogParts_MatchesSprintf(t *testing.T) {
	timestamp := time.Date(2026, 11, 1, 18, 30, 15, 123456789, time.FixedZone("CET", 3600))

	logParts := gosyslogformat.LogParts{
		"message":   "hello",
		"priority":  165,
		"version":   -1,
		"timestamp": timestamp,
		"now":       time.Now(),
		"bytes":     []byte("raw"),
		"other":     3.5,
	}

	record := parseLogParts(logParts)

	for key, value := range logParts {
		expected := fmt.Sprintf("%v", value)
		if bytes, ok := value.([]byte); ok {
			expected = string(bytes)
		}

		if record.Fields[key] != expected {
			t.Errorf("field %q: expected %q, got %q", key, expected, record.Fields[key])
		}
	}
}
