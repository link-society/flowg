package config

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestReadSystemConfig(t *testing.T) {
	ctx := t.Context()

	confOpts := DefaultOptions()
	confOpts.InMemory = true

	var confStorage Storage

	app := fxtest.New(
		t,
		NewStorage(confOpts),
		fx.Populate(&confStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	systemConfig, err := confStorage.ReadSystemConfig(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if systemConfig == nil {
		t.Fatal("systemConfig is nil")
	}

	if len(systemConfig.SyslogAllowedOrigins) != 0 {
		t.Fatal("AllowedOrigins is not empty")
	}

	systemConfig.SyslogAllowedOrigins = append(systemConfig.SyslogAllowedOrigins, "192.168.1.1")

	if err := confStorage.WriteSystemConfig(ctx, systemConfig); err != nil {
		t.Fatal(err)
	}

	systemConfig, err = confStorage.ReadSystemConfig(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if systemConfig == nil {
		t.Fatal("systemConfig is nil")
	}

	if len(systemConfig.SyslogAllowedOrigins) != 1 {
		t.Fatalf("len(AllowedOrigins) = %d != 1", len(systemConfig.SyslogAllowedOrigins))
	}

	if systemConfig.SyslogAllowedOrigins[0] != "192.168.1.1" {
		t.Fatal("systemConfig.AllowedOrigins[0] != 192.168.1.1")
	}
}
