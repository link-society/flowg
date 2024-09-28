package bootstrap_test

import (
	"context"
	"testing"

	"link-society.com/flowg/internal/storage/config"

	"link-society.com/flowg/internal/app/bootstrap"
)

func TestDefaultPipeline(t *testing.T) {
	ctx := context.Background()

	configStorage := config.NewStorage(config.OptInMemory(true))
	configStorage.Start()
	defer configStorage.Stop()

	err := configStorage.WaitStarted()
	if err != nil {
		t.Fatalf("unexpected error at startup: %v", err)
	}

	err = bootstrap.DefaultPipeline(ctx, configStorage)
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
