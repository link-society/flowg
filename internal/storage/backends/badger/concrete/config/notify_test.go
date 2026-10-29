package config_test

import (
	"testing"

	"log/slog"

	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"go.link-society.com/flowg/internal/models"

	"go.link-society.com/flowg/internal/engines/confignotify"
	confignotifyMocks "go.link-society.com/flowg/internal/engines/confignotify/mocks"

	badgerconfig "go.link-society.com/flowg/internal/storage/backends/badger/concrete/config"
	storage "go.link-society.com/flowg/internal/storage/interfaces"
)

// TestNotifyOnMutations verifies that every config mutation broadcasts the
// matching confignotify event.
func TestNotifyOnMutations(t *testing.T) {
	ctx := t.Context()

	notifier := confignotifyMocks.NewMockNotifier().(*confignotifyMocks.MockNotifier)
	notifier.On("Notify", mock.Anything, confignotify.Event{
		Kind:     confignotify.PipelineChanged,
		Pipeline: "test",
	}).Return(nil).Once()
	notifier.On("Notify", mock.Anything, confignotify.Event{
		Kind:     confignotify.PipelineDeleted,
		Pipeline: "test",
	}).Return(nil).Once()
	notifier.On("Notify", mock.Anything, confignotify.Event{
		Kind: confignotify.DependenciesChanged,
	}).Return(nil).Times(2)

	opts := badgerconfig.DefaultOptions()
	opts.InMemory = true

	var configStorage storage.ConfigStorage

	app := fxtest.New(
		t,
		fx.Supply(slog.New(slog.DiscardHandler)),
		fx.Supply(opts),
		badgerconfig.Module(),
		fx.Provide(func() confignotify.Notifier { return notifier }),
		fx.Populate(&configStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	if err := configStorage.WritePipeline(ctx, "test", &models.FlowGraphV2{
		MajorVersion: 2,
		MinorVersion: 1,
		Nodes:        []*models.FlowNodeV2{},
		Edges:        []*models.FlowEdgeV2{},
	}); err != nil {
		t.Fatalf("failed to write pipeline: %v", err)
	}

	if err := configStorage.DeletePipeline(ctx, "test"); err != nil {
		t.Fatalf("failed to delete pipeline: %v", err)
	}

	if err := configStorage.WriteTransformer(ctx, "test", "."); err != nil {
		t.Fatalf("failed to write transformer: %v", err)
	}

	if err := configStorage.DeleteTransformer(ctx, "test"); err != nil {
		t.Fatalf("failed to delete transformer: %v", err)
	}

	notifier.AssertExpectations(t)
}
