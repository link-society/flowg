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

	logM, err := notifier.Subscribe(ctx, "test")
	if err != nil {
		t.Fatalf("unexpected error while subscribing: %v", err)
	}

	logEntry := logstorage.NewLogEntry(map[string]string{})

	err = notifier.Notify(ctx, "test", "key", *logEntry)
	if err != nil {
		t.Fatalf("unexpected error while notifying: %v", err)
	}

	result, ok := <-logM.ReceiveC()
	if !ok {
		t.Fatalf("unexpected closed channel")
	}

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
