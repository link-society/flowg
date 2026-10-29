package bootstrap_test

import (
	"testing"

	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"go.link-society.com/flowg/pkg/models"

	"go.link-society.com/flowg/internal/storage/backends/badger/concrete/auth"
	storage "go.link-society.com/flowg/internal/storage/interfaces"

	"go.link-society.com/flowg/internal/storage/bootstrap"
)

func TestDefaultRolesAndUsers(t *testing.T) {
	ctx := t.Context()

	authOpts := auth.DefaultOptions()
	authOpts.InMemory = true

	var authStorage storage.AuthStorage

	app := fxtest.New(
		t,
		fx.Supply(slog.New(slog.DiscardHandler)),
		fx.Supply(authOpts),
		auth.Module(),
		fx.Populate(&authStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	err := bootstrap.DefaultRolesAndUsers(ctx, authStorage, bootstrap.BootstrapAuthOptions{
		InitialUser:     "root",
		InitialPassword: "root",
	})
	if err != nil {
		t.Fatalf("failed to bootstrap roles and users: %v", err)
	}

	roles, err := authStorage.ListRoles(ctx)
	if err != nil {
		t.Fatalf("failed to list roles: %v", err)
	}

	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}

	var adminRole *models.Role
	var viewerRole *models.Role

	for i := range roles {
		switch roles[i].Name {
		case "admin":
			adminRole = &roles[i]
		case "viewer":
			viewerRole = &roles[i]
		}
	}

	if adminRole == nil {
		t.Fatal("expected admin role to exist")
	}

	if viewerRole == nil {
		t.Fatal("expected viewer role to exist")
	}

	if !viewerRole.HasScope(models.SCOPE_READ_STREAMS) {
		t.Fatal("expected viewer role to have scope read_streams")
	}

	if viewerRole.HasScope(models.SCOPE_WRITE_STREAMS) {
		t.Fatal("expected viewer role to NOT have scope write_streams")
	}

	expected := []models.Scope{
		models.SCOPE_SEND_LOGS,
		models.SCOPE_WRITE_ACLS,
		models.SCOPE_WRITE_PIPELINES,
		models.SCOPE_WRITE_TRANSFORMERS,
		models.SCOPE_WRITE_STREAMS,
		models.SCOPE_WRITE_FORWARDERS,
	}

	for _, scope := range expected {
		if !adminRole.HasScope(scope) {
			t.Fatalf("expected admin role to have scope %s", scope)
		}
	}
}

func TestDefaultRolesAndUsers_DemoMode(t *testing.T) {
	ctx := t.Context()

	authOpts := auth.DefaultOptions()
	authOpts.InMemory = true

	var authStorage storage.AuthStorage

	app := fxtest.New(
		t,
		fx.Supply(slog.New(slog.DiscardHandler)),
		fx.Supply(authOpts),
		auth.Module(),
		fx.Populate(&authStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	err := bootstrap.DefaultRolesAndUsers(ctx, authStorage, bootstrap.BootstrapAuthOptions{
		InitialUser:     "root",
		InitialPassword: "root",
		DemoMode:        true,
	})
	if err != nil {
		t.Fatalf("failed to bootstrap roles and users: %v", err)
	}

	roles, err := authStorage.ListRoles(ctx)
	if err != nil {
		t.Fatalf("failed to list roles: %v", err)
	}

	var demoRole *models.Role
	for i := range roles {
		if roles[i].Name == "demo" {
			demoRole = &roles[i]
		}
	}

	if demoRole == nil {
		t.Fatal("expected demo role to exist")
	}

	if demoRole.HasScope(models.SCOPE_WRITE_ACLS) {
		t.Fatal("expected demo role to NOT have scope write_acls")
	}

	users, err := authStorage.ListUsers(ctx)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}

	var demoUser *models.User
	for i := range users {
		if users[i].Name == "demo" {
			demoUser = &users[i]
		}
	}

	if demoUser == nil {
		t.Fatal("expected demo user to exist")
	}
}
