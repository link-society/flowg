package bootstrap_test

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/internal/storage/config"

	"link-society.com/flowg/internal/app/bootstrap"
)

func TestDefaultPipeline(t *testing.T) {
	ctx := t.Context()

	var configStorage config.Storage

	configOpts := config.DefaultOptions()
	configOpts.InMemory = true

	app := fxtest.New(
		t,
		config.NewStorage(configOpts),
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
