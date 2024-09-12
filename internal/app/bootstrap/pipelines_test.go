package bootstrap_test

import (
	"testing"

	"link-society.com/flowg/internal/data/config"

	"link-society.com/flowg/internal/app/bootstrap"
)

func TestDefaultPipeline(t *testing.T) {
	opts := config.DefaultStorageOpts().WithInMemory(true)
	configStorage := config.NewStorage(opts)

	err := bootstrap.DefaultPipeline(configStorage)
	if err != nil {
		t.Fatalf("failed to create default pipeline: %v", err)
	}

	pipelineSys := config.NewPipelineSystem(configStorage)
	pipelines, err := pipelineSys.List()
	if err != nil {
		t.Fatalf("failed to list pipelines: %v", err)
	}

	if len(pipelines) != 1 {
		t.Fatalf("expected 1 pipeline, got %d", len(pipelines))
	}

	if pipelines[0] != "default" {
		t.Fatalf("expected pipeline name to be default, got %s", pipelines[0])
	}

	_, err = pipelineSys.Parse("default")
	if err != nil {
		t.Fatalf("failed to parse default pipeline: %v", err)
	}
}
