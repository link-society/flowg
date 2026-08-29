package config

import (
	"testing"

	"crypto/sha256"

	"github.com/stretchr/testify/mock"

	"link-society.com/flowg/internal/engines/confignotify"
	confignotifyMocks "link-society.com/flowg/internal/engines/confignotify/mocks"
)

type fakeSystemCache struct {
	invalidated int
}

func (f *fakeSystemCache) InvalidateSystemConfigCache() {
	f.invalidated++
}

// makeSnapshot hashes literal item values into the watcher's snapshot shape.
func makeSnapshot(items map[string]map[string]string) snapshot {
	result := make(snapshot, len(items))

	for itemType, byName := range items {
		hashes := make(map[string][sha256.Size]byte, len(byName))
		for name, value := range byName {
			hashes[name] = sha256.Sum256([]byte(value))
		}
		result[itemType] = hashes
	}

	return result
}

func newTestWatcher() (*watchWorker, *confignotifyMocks.MockNotifier, *fakeSystemCache) {
	notifier := confignotifyMocks.NewMockNotifier().(*confignotifyMocks.MockNotifier)
	cache := &fakeSystemCache{}

	return &watchWorker{
		notifier:    notifier,
		systemCache: cache,
	}, notifier, cache
}

func TestEmitDiff_SystemChangeInvalidatesCache(t *testing.T) {
	w, notifier, cache := newTestWatcher()

	previous := makeSnapshot(map[string]map[string]string{
		"system": {"config": "v1"},
	})
	next := makeSnapshot(map[string]map[string]string{
		"system": {"config": "v2"},
	})

	// no expectation set: any emitted event would fail the test
	w.emitDiff(t.Context(), previous, next)

	if cache.invalidated != 1 {
		t.Fatalf("expected 1 cache invalidation, got %d", cache.invalidated)
	}

	notifier.AssertExpectations(t)
}

func TestEmitDiff_PipelineChanges(t *testing.T) {
	w, notifier, cache := newTestWatcher()

	notifier.On("Notify", mock.Anything, confignotify.Event{
		Kind:     confignotify.PipelineChanged,
		Pipeline: "updated",
	}).Return(nil).Once()
	notifier.On("Notify", mock.Anything, confignotify.Event{
		Kind:     confignotify.PipelineChanged,
		Pipeline: "added",
	}).Return(nil).Once()
	notifier.On("Notify", mock.Anything, confignotify.Event{
		Kind:     confignotify.PipelineDeleted,
		Pipeline: "removed",
	}).Return(nil).Once()

	previous := makeSnapshot(map[string]map[string]string{
		"pipeline": {"updated": "v1", "removed": "v1", "untouched": "v1"},
	})
	next := makeSnapshot(map[string]map[string]string{
		"pipeline": {"updated": "v2", "added": "v1", "untouched": "v1"},
	})

	w.emitDiff(t.Context(), previous, next)

	if cache.invalidated != 0 {
		t.Fatalf("expected no cache invalidation, got %d", cache.invalidated)
	}

	notifier.AssertExpectations(t)
}

func TestEmitDiff_DependencyChangeCollapses(t *testing.T) {
	w, notifier, _ := newTestWatcher()

	// a single DependenciesChanged, even though a pipeline changed too
	notifier.On("Notify", mock.Anything, confignotify.Event{
		Kind: confignotify.DependenciesChanged,
	}).Return(nil).Once()

	previous := makeSnapshot(map[string]map[string]string{
		"transformer": {"parse": "v1"},
		"pipeline":    {"default": "v1"},
	})
	next := makeSnapshot(map[string]map[string]string{
		"transformer": {"parse": "v2"},
		"pipeline":    {"default": "v2"},
	})

	w.emitDiff(t.Context(), previous, next)

	notifier.AssertExpectations(t)
}
