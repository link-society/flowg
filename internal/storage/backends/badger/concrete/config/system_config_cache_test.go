package config_test

import (
	"testing"

	"bytes"
	"reflect"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/cmd/flowg-server/logging"
	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/backends/badger"
	badgerconfig "link-society.com/flowg/internal/storage/backends/badger/concrete/config"
	"link-society.com/flowg/internal/storage/databases/config"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

type adapterDeps struct {
	fx.In

	Adapter *badger.BadgerAdapter `name:"storage.config"`
}

// TestSystemConfigCacheInvalidation simulates two cluster nodes sharing the
// same storage: a write on node A stays invisible to node B's cache until node
// B invalidates it, which is what the FoundationDB watcher does on remote
// changes.
func TestSystemConfigCacheInvalidation(t *testing.T) {
	logging.Discard()

	ctx := t.Context()

	opts := badgerconfig.DefaultOptions()
	opts.InMemory = true

	var (
		nodeA storage.ConfigStorage
		d     adapterDeps
	)

	app := fxtest.New(
		t,
		badgerconfig.NewStorage(opts),
		fx.Populate(&nodeA),
		fx.Populate(&d),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	nodeB := config.NewStorage(d.Adapter, nil)

	v1 := &models.SystemConfiguration{SyslogAllowedOrigins: []string{"10.0.0.1"}}
	v2 := &models.SystemConfiguration{SyslogAllowedOrigins: []string{"10.0.0.2"}}

	if err := nodeA.WriteSystemConfig(ctx, v1); err != nil {
		t.Fatalf("failed to write system config: %v", err)
	}

	cached, err := nodeB.ReadSystemConfig(ctx)
	if err != nil {
		t.Fatalf("failed to read system config: %v", err)
	}
	if !reflect.DeepEqual(cached.SyslogAllowedOrigins, v1.SyslogAllowedOrigins) {
		t.Fatalf("unexpected system config: %v", cached)
	}

	if err := nodeA.WriteSystemConfig(ctx, v2); err != nil {
		t.Fatalf("failed to write system config: %v", err)
	}

	// node B still serves its cache
	stale, err := nodeB.ReadSystemConfig(ctx)
	if err != nil {
		t.Fatalf("failed to read system config: %v", err)
	}
	if !reflect.DeepEqual(stale.SyslogAllowedOrigins, v1.SyslogAllowedOrigins) {
		t.Fatalf("expected stale system config, got: %v", stale)
	}

	nodeB.InvalidateSystemConfigCache()

	fresh, err := nodeB.ReadSystemConfig(ctx)
	if err != nil {
		t.Fatalf("failed to read system config: %v", err)
	}
	if !reflect.DeepEqual(fresh.SyslogAllowedOrigins, v2.SyslogAllowedOrigins) {
		t.Fatalf("expected fresh system config, got: %v", fresh)
	}
}

// TestSystemConfigCacheInvalidatedOnRestore verifies that restoring a backup
// drops the cached system configuration.
func TestSystemConfigCacheInvalidatedOnRestore(t *testing.T) {
	logging.Discard()

	ctx := t.Context()

	v1 := &models.SystemConfiguration{SyslogAllowedOrigins: []string{"10.0.0.1"}}

	var dump bytes.Buffer

	// first instance: seed the config and take a backup
	{
		opts := badgerconfig.DefaultOptions()
		opts.InMemory = true

		var configStorage storage.ConfigStorage

		app := fxtest.New(
			t,
			badgerconfig.NewStorage(opts),
			fx.Populate(&configStorage),
			fx.NopLogger,
		)
		app.RequireStart()

		if err := configStorage.WriteSystemConfig(ctx, v1); err != nil {
			t.Fatalf("failed to write system config: %v", err)
		}
		if _, err := configStorage.Dump(ctx, &dump, 0); err != nil {
			t.Fatalf("failed to dump config: %v", err)
		}

		app.RequireStop()
	}

	// fresh instance: prime the cache with the empty config, then restore
	opts := badgerconfig.DefaultOptions()
	opts.InMemory = true

	var configStorage storage.ConfigStorage

	app := fxtest.New(
		t,
		badgerconfig.NewStorage(opts),
		fx.Populate(&configStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	primed, err := configStorage.ReadSystemConfig(ctx)
	if err != nil {
		t.Fatalf("failed to read system config: %v", err)
	}
	if len(primed.SyslogAllowedOrigins) != 0 {
		t.Fatalf("expected empty system config, got: %v", primed)
	}

	if err := configStorage.Load(ctx, &dump); err != nil {
		t.Fatalf("failed to restore config: %v", err)
	}

	restored, err := configStorage.ReadSystemConfig(ctx)
	if err != nil {
		t.Fatalf("failed to read system config: %v", err)
	}
	if !reflect.DeepEqual(restored.SyslogAllowedOrigins, v1.SyslogAllowedOrigins) {
		t.Fatalf("expected restored system config, got: %v", restored)
	}
}
