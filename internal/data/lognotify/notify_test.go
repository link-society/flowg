package lognotify_test

import (
	"reflect"
	"testing"

	"sync"
	"time"

	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/internal/data/lognotify"
)

func TestLogNotifier(t *testing.T) {
	notifier := lognotify.NewLogNotifier()
	notifier.Start()
	defer notifier.Stop()

	logDoneC := make(chan struct{})
	logC := notifier.Subscribe("test", logDoneC)

	logEntry := logstorage.NewLogEntry(map[string]string{})

	wgDoneC := make(chan struct{})
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

	go func() {
		wg.Wait()
		wgDoneC <- struct{}{}
		close(wgDoneC)
	}()

	select {
	case <-wgDoneC:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for log")
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

	logDoneC <- struct{}{}
	close(logDoneC)
}
