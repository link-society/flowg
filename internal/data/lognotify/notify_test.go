package lognotify_test

import (
	"context"
	"reflect"
	"testing"

	"time"

	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/internal/data/lognotify"
)

func TestLogNotifier(t *testing.T) {
	notifier := lognotify.NewLogNotifier()
	notifier.Start()
	defer notifier.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logM := notifier.Subscribe(ctx, "test")

	logEntry := logstorage.NewLogEntry(map[string]string{})

	notifier.Notify(ctx, "test", "key", *logEntry)
	result := <-logM.ReceiveC()

	if result.Stream != "test" {
		t.Fatalf("unexpected stream: %s", result.Stream)
	}

	if result.LogKey != "key" {
		t.Fatalf("unexpected log key: %s", result.LogKey)
	}

	if !reflect.DeepEqual(result.LogEntry, *logEntry) {
		t.Fatalf("unexpected log entry: %v", result.LogEntry)
	}
}
