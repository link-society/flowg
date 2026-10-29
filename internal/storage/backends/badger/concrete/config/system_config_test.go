package config_test

import (
	"testing"

	"log/slog"

	"sync"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"go.link-society.com/flowg/pkg/models"

	badgerconfig "go.link-society.com/flowg/internal/storage/backends/badger/concrete/config"
	storage "go.link-society.com/flowg/internal/storage/interfaces"
)

// TestReadSystemConfigConcurrent runs ReadSystemConfig against a concurrent
// writer. Under `go test -race` it flags the unsynchronized read of the cached
// system configuration unless the double-checked lock is made race-safe.
func TestReadSystemConfigConcurrent(t *testing.T) {
	ctx := t.Context()

	opts := badgerconfig.DefaultOptions()
	opts.InMemory = true

	var configStorage storage.ConfigStorage

	app := fxtest.New(
		t,
		fx.Supply(slog.New(slog.DiscardHandler)),
		fx.Supply(opts),
		badgerconfig.Module(),
		fx.Populate(&configStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	const (
		readers    = 8
		iterations = 100
	)

	var wg sync.WaitGroup

	for range readers {
		wg.Go(func() {
			for range iterations {
				if _, err := configStorage.ReadSystemConfig(ctx); err != nil {
					t.Errorf("failed to read system config: %v", err)
					return
				}
			}
		})
	}

	wg.Go(func() {
		for j := 0; j < iterations; j++ {
			if err := configStorage.WriteSystemConfig(ctx, &models.SystemConfiguration{}); err != nil {
				t.Errorf("failed to write system config: %v", err)
				return
			}
		}
	})

	wg.Wait()
}
