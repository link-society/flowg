package forwarders

import (
	"testing"
	"time"

	"link-society.com/flowg/internal/models"
)

func TestAwsCloudWatchTimestamp(t *testing.T) {
	timestamp := time.Date(2026, time.August, 5, 15, 0, 0, 123_000_000, time.UTC)
	record := &models.LogRecord{Timestamp: timestamp}

	if got := awsCloudWatchTimestamp(record); got != timestamp.UnixMilli() {
		t.Fatalf("expected timestamp %d, got %d", timestamp.UnixMilli(), got)
	}
}
