package pipelines

import (
	"testing"

	"errors"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/engines/confignotify"
	"link-society.com/flowg/internal/engines/lognotify"

	badgerconfig "link-society.com/flowg/internal/storage/backends/badger/concrete/config"
	badgerlog "link-society.com/flowg/internal/storage/backends/badger/concrete/log"
)

const rawTestPipeline = `{
	"version": 2,
	"nodes": [
		{
			"id": "source",
			"type": "source",
			"data": {
				"type": "direct"
			}
		}
	],
	"edges": []
}`

func testRunnerModules(configOpts badgerconfig.Options, logOpts badgerlog.Options) []fx.Option {
	return []fx.Option{
		confignotify.NewNotifier(),
		badgerconfig.NewStorage(configOpts),
		badgerlog.NewStorage(logOpts),
		lognotify.NewLogNotifier(),
		NewRunner(),
		fx.NopLogger,
	}
}

func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("condition not met before timeout")
}

func TestRunner_StorageDrivenLifecycle(t *testing.T) {
	configOpts := badgerconfig.DefaultOptions()
	configOpts.InMemory = true
	logOpts := badgerlog.DefaultOptions()
	logOpts.InMemory = true

	var (
		runner        Runner
		configStorage storage.ConfigStorage
	)

	app := fxtest.New(
		t,
		fx.Options(testRunnerModules(configOpts, logOpts)...),
		fx.Populate(&runner, &configStorage),
	)
	app.RequireStart()
	defer app.RequireStop()

	ctx := t.Context()
	record := models.NewLogRecord(map[string]string{"message": "test"})

	process := func() error {
		return runner.Process(ctx, "test", DIRECT_ENTRYPOINT, record)
	}

	var notFound *PipelineNotFoundError
	if err := process(); !errors.As(err, &notFound) {
		t.Fatalf("expected PipelineNotFoundError before creation, got: %v", err)
	}

	// creating the pipeline in storage starts it
	if err := configStorage.WriteRawPipeline(ctx, "test", rawTestPipeline); err != nil {
		t.Fatalf("failed to write pipeline: %v", err)
	}
	waitFor(t, 5*time.Second, func() bool { return process() == nil })

	// deleting it from storage terminates it
	if err := configStorage.DeletePipeline(ctx, "test"); err != nil {
		t.Fatalf("failed to delete pipeline: %v", err)
	}
	waitFor(t, 5*time.Second, func() bool {
		return errors.As(process(), &notFound)
	})
}

func TestRunner_StartsPersistedPipelinesOnBoot(t *testing.T) {
	dir := t.TempDir()

	configOpts := badgerconfig.DefaultOptions()
	configOpts.Directory = dir
	logOpts := badgerlog.DefaultOptions()
	logOpts.InMemory = true

	// seed the pipeline with a first app instance
	{
		var configStorage storage.ConfigStorage

		app := fxtest.New(
			t,
			fx.Options(testRunnerModules(configOpts, logOpts)...),
			fx.Populate(&configStorage),
		)
		app.RequireStart()

		if err := configStorage.WriteRawPipeline(t.Context(), "test", rawTestPipeline); err != nil {
			t.Fatalf("failed to write pipeline: %v", err)
		}

		app.RequireStop()
	}

	// a fresh instance on the same storage starts it eagerly
	var runner Runner

	app := fxtest.New(
		t,
		fx.Options(testRunnerModules(configOpts, logOpts)...),
		fx.Populate(&runner),
	)
	app.RequireStart()
	defer app.RequireStop()

	record := models.NewLogRecord(map[string]string{"message": "test"})
	if err := runner.Process(t.Context(), "test", DIRECT_ENTRYPOINT, record); err != nil {
		t.Fatalf("expected pipeline to be running after boot, got: %v", err)
	}
}
