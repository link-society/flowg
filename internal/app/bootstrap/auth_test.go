package bootstrap_test

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/auth"

	"link-society.com/flowg/internal/app/bootstrap"
	"link-society.com/flowg/internal/app/logging"
)

func TestDefaultRolesAndUsers(t *testing.T) {
	logging.Discard()

	ctx := t.Context()

	authOpts := auth.DefaultOptions()
	authOpts.InMemory = true

	var authStorage auth.Storage

	app := fxtest.New(
		t,
		auth.NewStorage(authOpts),
		fx.Populate(&authStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	defer app.RequireStop()

	err := bootstrap.DefaultRolesAndUsers(ctx, authStorage, bootstrap.BootstrapOptions{
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

	if len(roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(roles))
	}

	if roles[0].Name != "admin" {
		t.Fatalf("expected role name to be admin, got %s", roles[0].Name)
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
		if !roles[0].HasScope(scope) {
			t.Fatalf("expected role to have scope %s", scope)
		}
	}
}
