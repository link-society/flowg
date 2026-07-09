package bootstrap_test

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/backends/badger/concrete/auth"
	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/cmd/flowg-server/logging"
	"link-society.com/flowg/internal/storage/bootstrap"
)

func TestDefaultRolesAndUsers(t *testing.T) {
	logging.Discard()

	ctx := t.Context()

	authOpts := auth.DefaultOptions()
	authOpts.InMemory = true

	var authStorage storage.AuthStorage

	app := fxtest.New(
		t,
		auth.NewStorage(authOpts),
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
