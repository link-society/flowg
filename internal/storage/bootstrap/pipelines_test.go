package bootstrap_test

import (
	"testing"

	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"go.link-society.com/flowg/internal/storage/backends/badger/concrete/config"
	storage "go.link-society.com/flowg/internal/storage/interfaces"

	"go.link-society.com/flowg/internal/storage/bootstrap"
)

func TestDefaultPipeline(t *testing.T) {
	ctx := t.Context()

	var configStorage storage.ConfigStorage

	configOpts := config.DefaultOptions()
	configOpts.InMemory = true

	app := fxtest.New(
		t,
		fx.Supply(slog.New(slog.DiscardHandler)),
		fx.Supply(configOpts),
		config.Module(),
		fx.Populate(&configStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	err := bootstrap.DefaultPipeline(ctx, configStorage)
	if err != nil {
		t.Fatalf("failed to bootstrap default pipeline: %v", err)
	}

	pipelines, err := configStorage.ListPipelines(ctx)
	if err != nil {
		t.Fatalf("failed to list pipelines: %v", err)
	}

	if len(pipelines) != 1 {
		t.Fatalf("expected 1 pipeline, got %d", len(pipelines))
	}

	if pipelines[0] != "default" {
		t.Fatalf("expected pipeline name to be default, got %s", pipelines[0])
	}

	_, err = configStorage.ReadPipeline(ctx, "default")
	if err != nil {
		t.Fatalf("failed to parse default pipeline: %v", err)
	}
}
