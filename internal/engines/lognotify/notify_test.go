package lognotify_test

import (
	"fmt"
	"reflect"
	"testing"

	"context"
	"time"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/engines/lognotify"
)

func TestLogNotifier(t *testing.T) {
	fmt.Printf("lognotify: starting log notifier")
	notifier := lognotify.NewLogNotifier()
	notifier.Start()
	defer notifier.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logM, err := notifier.Subscribe(ctx, "test")
	if err != nil {
		t.Fatalf("unexpected error while subscribing: %v", err)
	}

	logRecord := models.NewLogRecord(map[string]string{})

	err = notifier.Notify(ctx, "test", "key", *logRecord)
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

	if !reflect.DeepEqual(result.LogRecord, *logRecord) {
		t.Fatalf("unexpected log record: %v", result.LogRecord)
	}
}
