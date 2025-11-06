package lognotify_test

import (
	"reflect"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/engines/lognotify"
)

func TestLogNotifier(t *testing.T) {
	var notifier lognotify.LogNotifier

	app := fxtest.New(
		t,
		lognotify.NewLogNotifier(),
		fx.Populate(&notifier),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	logM, err := notifier.Subscribe(t.Context(), "test")
	if err != nil {
		t.Fatalf("unexpected error while subscribing: %v", err)
	}

	logRecord := models.NewLogRecord(map[string]string{})

	err = notifier.Notify(t.Context(), "test", "key", *logRecord)
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
