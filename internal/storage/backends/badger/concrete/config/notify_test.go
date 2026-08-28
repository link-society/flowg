package config_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/cmd/flowg-server/logging"
	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/engines/confignotify"
	confignotifyMocks "link-society.com/flowg/internal/engines/confignotify/mocks"

	badgerconfig "link-society.com/flowg/internal/storage/backends/badger/concrete/config"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// TestNotifyOnMutations verifies that every config mutation broadcasts the
// matching confignotify event.
func TestNotifyOnMutations(t *testing.T) {
	logging.Discard()

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
		badgerconfig.NewStorage(opts),
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
