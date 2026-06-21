package bootstrap_test

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/internal/app/bootstrap"
	"link-society.com/flowg/internal/app/logging"

	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/utils/hlc"
)

func TestDefaultSystemConfig(t *testing.T) {
	logging.Discard()

	ctx := t.Context()

	confOpts := config.DefaultOptions()
	confOpts.InMemory = true

	var confStorage config.Storage

	app := fxtest.New(
		t,
		fx.Provide(func() *hlc.Clock { return hlc.NewClock("test") }),
		config.NewStorage(confOpts),
		fx.Populate(&confStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	err := bootstrap.DefaultSystemConfig(ctx, confStorage, bootstrap.BootstrapSystemOptions{
		InitialSyslogAllowedOrigins: []string{"127.0.0.1"},
	})
	if err != nil {
		t.Fatalf("failed to bootstrap default system config: %v", err)
	}

	systemConfig, err := confStorage.ReadSystemConfig(ctx)
	if err != nil {
		t.Fatalf("failed to read system config: %v", err)
	}

	if len(systemConfig.SyslogAllowedOrigins) != 1 {
		t.Fatalf("expected 1 allowed origin, got %d", len(systemConfig.SyslogAllowedOrigins))
	}

	if systemConfig.SyslogAllowedOrigins[0] != "127.0.0.1" {
		t.Fatalf("expected allowed origin to be 127.0.0.1, got %s", systemConfig.SyslogAllowedOrigins[0])
	}
}
