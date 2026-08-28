package confignotify_test

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/internal/engines/confignotify"
)

func TestNotifier(t *testing.T) {
	var notifier confignotify.Notifier

	app := fxtest.New(
		t,
		confignotify.NewNotifier(),
		fx.Populate(&notifier),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	eventM, err := notifier.Subscribe(t.Context())
	if err != nil {
		t.Fatalf("unexpected error while subscribing: %v", err)
	}

	event := confignotify.Event{
		Kind:     confignotify.PipelineChanged,
		Pipeline: "test",
	}

	if err := notifier.Notify(t.Context(), event); err != nil {
		t.Fatalf("unexpected error while notifying: %v", err)
	}

	result, ok := <-eventM.ReceiveC()
	if !ok {
		t.Fatalf("unexpected closed channel")
	}

	if result != event {
		t.Fatalf("unexpected event: %v", result)
	}
}

func TestNotifier_MultipleSubscribers(t *testing.T) {
	var notifier confignotify.Notifier

	app := fxtest.New(
		t,
		confignotify.NewNotifier(),
		fx.Populate(&notifier),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	first, err := notifier.Subscribe(t.Context())
	if err != nil {
		t.Fatalf("unexpected error while subscribing: %v", err)
	}
	second, err := notifier.Subscribe(t.Context())
	if err != nil {
		t.Fatalf("unexpected error while subscribing: %v", err)
	}

	event := confignotify.Event{Kind: confignotify.DependenciesChanged}

	if err := notifier.Notify(t.Context(), event); err != nil {
		t.Fatalf("unexpected error while notifying: %v", err)
	}

	for _, sub := range []interface{ ReceiveC() <-chan confignotify.Event }{first, second} {
		result, ok := <-sub.ReceiveC()
		if !ok {
			t.Fatalf("unexpected closed channel")
		}
		if result != event {
			t.Fatalf("unexpected event: %v", result)
		}
	}
}
