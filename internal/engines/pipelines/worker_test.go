package pipelines

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"context"
	"errors"

	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/models"

	lognotifyMocks "link-society.com/flowg/internal/engines/lognotify/mocks"
	"link-society.com/flowg/internal/storage/generic/kv"
	"link-society.com/flowg/internal/storage/mocks"

	appMetrics "link-society.com/flowg/internal/app/metrics"
)

func TestMain(m *testing.M) {
	appMetrics.Setup()
	m.Run()
}

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

// routedFlowGraph returns a flow graph whose direct source feeds a router node,
// so that processing exercises the log storage.
func routedFlowGraph() *models.FlowGraphV2 {
	return &models.FlowGraphV2{
		MajorVersion: 2,
		MinorVersion: 1,
		Nodes: []*models.FlowNodeV2{
			{
				ID:   "source",
				Type: "source",
				Data: map[string]string{"type": "direct"},
			},
			{
				ID:   "router",
				Type: "router",
				Data: map[string]string{"stream": "default"},
			},
		},
		Edges: []*models.FlowEdgeV2{
			{ID: "e1", Source: "source", Target: "router"},
		},
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

func TestClosePipeline_DrainsBeforeClose(t *testing.T) {
	configStorage := mocks.NewMockConfigStorage().(*mocks.MockConfigStorage)
	configStorage.On("ReadPipeline", mock.Anything, "test").
		Return(routedFlowGraph(), nil)

	entered := make(chan struct{})
	unblock := make(chan struct{})

	logStorage := mocks.NewMockLogStorage().(*mocks.MockLogStorage)
	logStorage.On("Ingest", mock.Anything, "default", mock.Anything).
		Run(func(args mock.Arguments) {
			close(entered)
			<-unblock
		}).
		Return(kv.Key{"k"}, nil)

	logNotifier := lognotifyMocks.NewMockNotifier().(*lognotifyMocks.MockNotifier)
	logNotifier.On("Notify", mock.Anything, "default", mock.Anything, mock.Anything).
		Return(nil)

	w := newTestWorker(configStorage)
	w.logStorage = logStorage
	w.logNotifier = logNotifier

	if err := w.buildPipeline(actor.ContextStarted(), "test"); err != nil {
		t.Fatalf("failed to build pipeline: %v", err)
	}

	// in-flight record, blocked inside the router node
	processDone := make(chan error, 1)
	pipeline, release, err := w.acquirePipeline("test")
	if err != nil {
		t.Fatalf("failed to acquire pipeline: %v", err)
	}
	go func() {
		defer release()
		ctx := context.WithValue(context.Background(), workerKey, w)
		processDone <- pipeline.Process(ctx, DIRECT_ENTRYPOINT, &models.LogRecord{
			Timestamp: time.Now(),
			Fields:    map[string]string{"message": "test"},
		})
	}()

	<-entered

	closeDone, err := w.closePipeline("test")
	if err != nil {
		t.Fatalf("failed to close pipeline: %v", err)
	}

	// new uses are rejected while the record is still in flight
	if _, _, err := w.acquirePipeline("test"); err == nil {
		t.Fatal("expected acquire to fail after close")
	}

	// the close must not complete while the record is in flight
	select {
	case <-closeDone:
		t.Fatal("pipeline closed while a record was in flight")
	case <-time.After(50 * time.Millisecond):
	}

	close(unblock)

	select {
	case err := <-closeDone:
		if err != nil {
			t.Fatalf("failed to close pipeline: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("close did not complete after drain")
	}

	if err := <-processDone; err != nil {
		t.Fatalf("in-flight record failed: %v", err)
	}
}

func TestBuildPipeline_SwapRetiresPrevious(t *testing.T) {
	configStorage := mocks.NewMockConfigStorage().(*mocks.MockConfigStorage)
	configStorage.On("ReadPipeline", mock.Anything, "test").
		Return(minimalFlowGraph(), nil)

	w := newTestWorker(configStorage)

	if err := w.buildPipeline(actor.ContextStarted(), "test"); err != nil {
		t.Fatalf("failed to build pipeline: %v", err)
	}
	first, releaseFirst, err := w.acquirePipeline("test")
	if err != nil {
		t.Fatalf("failed to acquire first build: %v", err)
	}

	if err := w.buildPipeline(actor.ContextStarted(), "test"); err != nil {
		t.Fatalf("failed to rebuild pipeline: %v", err)
	}

	second, releaseSecond, err := w.acquirePipeline("test")
	if err != nil {
		t.Fatalf("failed to acquire second build: %v", err)
	}
	defer releaseSecond()

	if first == second {
		t.Fatal("expected rebuild to produce a new build")
	}

	// the retired build refuses new uses but stays valid for in-flight ones
	if first.acquire() {
		t.Fatal("expected first build to be retired")
	}
	releaseFirst()
}
