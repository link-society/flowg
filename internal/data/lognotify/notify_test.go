package lognotify_test

import (
	"reflect"
	"sync"
	"testing"

	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/internal/data/lognotify"
)

func TestLogNotifier(t *testing.T) {
	notifier := lognotify.NewLogNotifier()
	notifier.Start()
	defer notifier.Stop()

	doneC := make(chan struct{})
	logC := notifier.Subscribe("test", doneC)

	logEntry := logstorage.NewLogEntry(map[string]string{})

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		notifier.Notify("test", "key", *logEntry)
	}()

	var result lognotify.LogMessage
	go func() {
		defer wg.Done()
		result = <-logC
	}()

	wg.Wait()

	if result.Stream != "test" {
		t.Fatalf("unexpected stream: %s", result.Stream)
	}

	if result.LogKey != "key" {
		t.Fatalf("unexpected log key: %s", result.LogKey)
	}

	if !reflect.DeepEqual(result.LogEntry, *logEntry) {
		t.Fatalf("unexpected log entry: %v", result.LogEntry)
	}

	doneC <- struct{}{}
	close(doneC)
}
