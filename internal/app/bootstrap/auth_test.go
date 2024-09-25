package bootstrap_test

import (
	"testing"

	"link-society.com/flowg/internal/data/auth"

	"link-society.com/flowg/internal/app/bootstrap"
	"link-society.com/flowg/internal/app/logging"
)

func TestDefaultRolesAndUsers(t *testing.T) {
	logging.Discard()

	opts := auth.DefaultDatabaseOpts().WithInMemory(true)
	authDb := auth.NewDatabase(opts)
	err := authDb.Open()
	if err != nil {
		t.Fatalf("failed to create auth database: %v", err)
	}
	defer authDb.Close()

	err = bootstrap.DefaultRolesAndUsers(authDb)
	if err != nil {
		t.Fatalf("failed to bootstrap roles and users: %v", err)
	}

	roleSys := auth.NewRoleSystem(authDb)

	roles, err := roleSys.ListRoles()
	if err != nil {
		t.Fatalf("failed to list roles: %v", err)
	}

	if len(roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(roles))
	}

	if roles[0].Name != "admin" {
		t.Fatalf("expected role name to be admin, got %s", roles[0].Name)
	}

	expected := []auth.Scope{
		auth.SCOPE_SEND_LOGS,
		auth.SCOPE_WRITE_ACLS,
		auth.SCOPE_WRITE_PIPELINES,
		auth.SCOPE_WRITE_TRANSFORMERS,
		auth.SCOPE_WRITE_STREAMS,
		auth.SCOPE_WRITE_ALERTS,
	}

	for _, scope := range expected {
		if !roles[0].HasScope(scope) {
			t.Fatalf("expected role to have scope %s", scope)
		}
	}
}
