package pipelines

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/mocks"
)

// minimalFlowGraph returns a flow graph with a single direct source node,
// enough to compile without touching transformers or forwarders.
func minimalFlowGraph() *models.FlowGraphV2 {
	return &models.FlowGraphV2{
		MajorVersion: 2,
		MinorVersion: 1,
		Nodes: []*models.FlowNodeV2{
			{
				ID:   "source",
				Type: "source",
				Data: map[string]string{"type": "direct"},
			},
		},
		Edges: []*models.FlowEdgeV2{},
	}
}

func newTestWorker(configStorage *mocks.MockConfigStorage) *worker {
	return &worker{
		configStorage: configStorage,
		cache:         make(map[string]*Pipeline),
	}
}

func TestBuildPipeline_UnknownPipeline(t *testing.T) {
	configStorage := mocks.NewMockConfigStorage().(*mocks.MockConfigStorage)
	configStorage.On("ReadPipeline", mock.Anything, "missing").
		Return((*models.FlowGraphV2)(nil), nil)

	w := newTestWorker(configStorage)

	err := w.buildPipeline(actor.ContextStarted(), "missing")

	var notFound *PipelineNotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected PipelineNotFoundError, got: %v", err)
	}

	if len(w.cache) != 0 {
		t.Fatalf("expected empty cache, got %d entries", len(w.cache))
	}
}

func TestBuildPipeline_AddsToCache(t *testing.T) {
	configStorage := mocks.NewMockConfigStorage().(*mocks.MockConfigStorage)
	configStorage.On("ReadPipeline", mock.Anything, "test").
		Return(minimalFlowGraph(), nil)

	w := newTestWorker(configStorage)

	if err := w.buildPipeline(actor.ContextStarted(), "test"); err != nil {
		t.Fatalf("failed to build pipeline: %v", err)
	}

	if _, exists := w.cache["test"]; !exists {
		t.Fatal("expected pipeline in cache")
	}
}
