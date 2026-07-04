//go:build integration_fdb

package auth_test

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"

	fdb_auth "link-society.com/flowg/internal/storage/backends/foundationdb/concrete/auth"
)

func connectString() string {
	return "docker:docker@127.0.0.1:4500"
}

func newAuthStorage(t *testing.T) storage.AuthStorage {
	t.Helper()

	opts := fdb_auth.DefaultOptions()
	opts.ConnectionString = connectString()

	var authStorage storage.AuthStorage

	app := fxtest.New(
		t,
		fdb_auth.NewStorage(opts),
		fx.Populate(&authStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	t.Cleanup(app.RequireStop)

	return authStorage
}

// TestAuth_All runs all auth tests against a single FDB connection to avoid
// triggering cgo finalizer crashes from rapid connection create/destroy.
func TestAuth_All(t *testing.T) {
	ctx := t.Context()
	auth := newAuthStorage(t)

	// ========================================================================
	//  Roles
	// ========================================================================
	t.Log("=== Roles ===")

	// --- ListRoles empty ---
	roles, err := auth.ListRoles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 0 {
		t.Fatalf("expected 0 roles initially, got %d", len(roles))
	}

	// --- Create role ---
	role := models.Role{
		Name:   "test-role",
		Scopes: []models.Scope{models.SCOPE_READ_STREAMS, models.SCOPE_SEND_LOGS},
	}
	if err := auth.SaveRole(ctx, role); err != nil {
		t.Fatal(err)
	}

	// --- Fetch role ---
	fetched, err := auth.FetchRole(ctx, "test-role")
	if err != nil {
		t.Fatal(err)
	}
	if fetched == nil {
		t.Fatal("FetchRole: expected role, got nil")
	}
	if fetched.Name != "test-role" {
		t.Fatalf("FetchRole: expected 'test-role', got '%s'", fetched.Name)
	}
	if len(fetched.Scopes) != 2 {
		t.Fatalf("FetchRole: expected 2 scopes, got %d", len(fetched.Scopes))
	}

	// --- List one role ---
	roles, err = auth.ListRoles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 1 {
		t.Fatalf("ListRoles: expected 1 role, got %d", len(roles))
	}

	// --- Fetch non-existent role ---
	missing, err := auth.FetchRole(ctx, "no-such-role")
	if err != nil {
		t.Fatal(err)
	}
	if missing != nil {
		t.Fatal("FetchRole non-existent: expected nil")
	}

	// --- Update role scopes ---
	role.Scopes = []models.Scope{models.SCOPE_READ_STREAMS, models.SCOPE_READ_PIPELINES, models.SCOPE_WRITE_FORWARDERS}
	if err := auth.SaveRole(ctx, role); err != nil {
		t.Fatal(err)
	}
	fetched, err = auth.FetchRole(ctx, "test-role")
	if err != nil {
		t.Fatal(err)
	}
	if len(fetched.Scopes) != 3 {
		t.Fatalf("UpdateRole: expected 3 scopes, got %d", len(fetched.Scopes))
	}

	// --- Duplicate save ---
	if err := auth.SaveRole(ctx, role); err != nil {
		t.Fatal(err)
	}
	roles, err = auth.ListRoles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 1 {
		t.Fatalf("DuplicateSave: expected 1 role, got %d", len(roles))
	}

	// --- Delete non-existent ---
	if err := auth.DeleteRole(ctx, "non-existent"); err != nil {
		t.Fatal(err)
	}

	// --- Delete role ---
	if err := auth.DeleteRole(ctx, "test-role"); err != nil {
		t.Fatal(err)
	}
	fetched, err = auth.FetchRole(ctx, "test-role")
	if err != nil {
		t.Fatal(err)
	}
	if fetched != nil {
		t.Fatal("DeleteRole: expected nil after delete")
	}
	roles, err = auth.ListRoles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 0 {
		t.Fatalf("DeleteRole: expected 0 roles, got %d", len(roles))
	}

	// ========================================================================
	//  Users
	// ========================================================================
	t.Log("=== Users ===")

	// --- ListUsers empty ---
	users, err := auth.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 0 {
		t.Fatalf("ListUsers: expected 0 users, got %d", len(users))
	}

	// Create roles for user assignment
	auth.SaveRole(ctx, models.Role{Name: "viewer", Scopes: []models.Scope{models.SCOPE_READ_STREAMS, models.SCOPE_SEND_LOGS}})
	auth.SaveRole(ctx, models.Role{Name: "admin", Scopes: []models.Scope{models.SCOPE_WRITE_PIPELINES, models.SCOPE_WRITE_ACLS}})

	// --- Create user ---
	user := models.User{Name: "testuser", Roles: []string{"viewer"}}
	if err := auth.SaveUser(ctx, user, "secret123"); err != nil {
		t.Fatal(err)
	}

	// --- Fetch user ---
	fetchedUser, err := auth.FetchUser(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if fetchedUser == nil {
		t.Fatal("FetchUser: expected user, got nil")
	}
	if fetchedUser.Name != "testuser" {
		t.Fatalf("FetchUser: expected 'testuser', got '%s'", fetchedUser.Name)
	}
	if len(fetchedUser.Roles) != 1 || fetchedUser.Roles[0] != "viewer" {
		t.Fatalf("FetchUser: expected [viewer], got %v", fetchedUser.Roles)
	}

	// --- Fetch non-existent user ---
	missingUser, err := auth.FetchUser(ctx, "nobody")
	if err != nil {
		t.Fatal(err)
	}
	if missingUser != nil {
		t.Fatal("FetchUser non-existent: expected nil")
	}

	// --- List users ---
	users, err = auth.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatalf("ListUsers: expected 1 user, got %d", len(users))
	}

	// --- Patch user roles (add) ---
	user.Roles = []string{"viewer", "admin"}
	if err := auth.PatchUserRoles(ctx, user); err != nil {
		t.Fatal(err)
	}
	fetchedUser, err = auth.FetchUser(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if len(fetchedUser.Roles) != 2 {
		t.Fatalf("PatchRoles: expected 2 roles, got %d", len(fetchedUser.Roles))
	}

	// --- Patch user roles (remove) ---
	user.Roles = []string{"admin"}
	if err := auth.PatchUserRoles(ctx, user); err != nil {
		t.Fatal(err)
	}
	fetchedUser, err = auth.FetchUser(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if len(fetchedUser.Roles) != 1 || fetchedUser.Roles[0] != "admin" {
		t.Fatalf("PatchRoles remove: expected [admin], got %v", fetchedUser.Roles)
	}

	// ========================================================================
	//  Passwords
	// ========================================================================
	t.Log("=== Passwords ===")

	// --- Correct password ---
	ok, err := auth.VerifyUserPassword(ctx, "testuser", "secret123")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("VerifyPassword: expected true")
	}

	// --- Wrong password ---
	ok, err = auth.VerifyUserPassword(ctx, "testuser", "wrong-password")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("VerifyPassword: expected false for wrong password")
	}

	// --- Empty password ---
	ok, err = auth.VerifyUserPassword(ctx, "testuser", "")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("VerifyPassword: expected false for empty password")
	}

	// --- Non-existent user ---
	_, err = auth.VerifyUserPassword(ctx, "ghost", "any")
	if err == nil {
		t.Fatal("VerifyPassword: expected error for non-existent user")
	}

	// ========================================================================
	//  Permissions
	// ========================================================================
	t.Log("=== Permissions ===")

	// User has "admin" role with write_pipelines + write_acls

	ok, err = auth.VerifyUserPermission(ctx, "testuser", models.SCOPE_WRITE_PIPELINES)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Permission: expected write_pipelines")
	}

	// Write implies read
	ok, err = auth.VerifyUserPermission(ctx, "testuser", models.SCOPE_READ_PIPELINES)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Permission: expected write_pipelines to imply read_pipelines")
	}

	// No unrelated scope
	ok, err = auth.VerifyUserPermission(ctx, "testuser", models.SCOPE_READ_STREAMS)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("Permission: expected false for read_streams")
	}

	// Non-existent user
	ok, err = auth.VerifyUserPermission(ctx, "nobody", models.SCOPE_READ_STREAMS)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("Permission: expected false for non-existent user")
	}

	// ========================================================================
	//  ListUserScopes
	// ========================================================================
	t.Log("=== ListUserScopes ===")

	scopes, err := auth.ListUserScopes(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	expectedScopes := map[models.Scope]bool{
		models.SCOPE_WRITE_PIPELINES: true,
		models.SCOPE_READ_PIPELINES:  true,
		models.SCOPE_WRITE_ACLS:      true,
		models.SCOPE_READ_ACLS:       true,
	}
	for _, s := range scopes {
		if !expectedScopes[s] {
			t.Fatalf("ListUserScopes: unexpected scope '%s'", s)
		}
		expectedScopes[s] = false
	}
	for s, seen := range expectedScopes {
		if seen {
			t.Fatalf("ListUserScopes: missing scope '%s'", s)
		}
	}

	emptyScopes, err := auth.ListUserScopes(ctx, "ghost")
	if err != nil {
		t.Fatal(err)
	}
	if len(emptyScopes) != 0 {
		t.Fatalf("ListUserScopes ghost: expected 0, got %d", len(emptyScopes))
	}

	// ========================================================================
	//  Tokens
	// ========================================================================
	t.Log("=== Tokens ===")

	// Restore both roles for token test
	user.Roles = []string{"viewer", "admin"}
	auth.PatchUserRoles(ctx, user)

	// --- Create token ---
	token, tokenUUID, err := auth.CreateToken(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("CreateToken: expected non-empty token")
	}
	if tokenUUID == "" {
		t.Fatal("CreateToken: expected non-empty UUID")
	}

	// --- List tokens ---
	tokens, err := auth.ListTokens(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 1 || tokens[0] != tokenUUID {
		t.Fatalf("ListTokens: expected [%s], got %v", tokenUUID, tokens)
	}

	// --- Verify token ---
	verifiedUser, err := auth.VerifyToken(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	if verifiedUser == nil {
		t.Fatal("VerifyToken: expected user, got nil")
	}
	if verifiedUser.Name != "testuser" {
		t.Fatalf("VerifyToken: expected 'testuser', got '%s'", verifiedUser.Name)
	}

	// --- Invalid token ---
	badUser, err := auth.VerifyToken(ctx, "pat_invalid_token_12345")
	if err != nil {
		t.Fatal(err)
	}
	if badUser != nil {
		t.Fatal("VerifyToken: expected nil for invalid token")
	}

	// --- Multiple tokens ---
	extraToken, extraUUID, err := auth.CreateToken(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = auth.CreateToken(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	tokens, err = auth.ListTokens(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 3 {
		t.Fatalf("MultipleTokens: expected 3, got %d", len(tokens))
	}

	// --- Delete one token ---
	if err := auth.DeleteToken(ctx, "testuser", extraUUID); err != nil {
		t.Fatal(err)
	}
	tokens, err = auth.ListTokens(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 2 {
		t.Fatalf("DeleteToken: expected 2, got %d", len(tokens))
	}

	badUser, err = auth.VerifyToken(ctx, extraToken)
	if err != nil {
		t.Fatal(err)
	}
	if badUser != nil {
		t.Fatal("VerifyToken: expected nil for deleted token")
	}

	// Original token still works
	verifiedUser, err = auth.VerifyToken(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	if verifiedUser == nil {
		t.Fatal("VerifyToken: original token should still work")
	}

	// ========================================================================
	//  User deletion cascades to tokens
	// ========================================================================
	t.Log("=== User deletion cascades ===")

	if err := auth.DeleteUser(ctx, "testuser"); err != nil {
		t.Fatal(err)
	}

	fetchedUser, err = auth.FetchUser(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if fetchedUser != nil {
		t.Fatal("DeleteUser: expected nil")
	}

	badUser, err = auth.VerifyToken(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	if badUser != nil {
		t.Fatal("DeleteUser: token should be invalid after user deletion")
	}

	users, err = auth.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 0 {
		t.Fatalf("DeleteUser: expected 0 users, got %d", len(users))
	}

	// ========================================================================
	//  PatchUserRoles preserves password
	// ========================================================================
	t.Log("=== Patch preserves password ===")

	auth.SaveUser(ctx, models.User{Name: "pw-preserve", Roles: []string{"viewer"}}, "initial-pw")
	auth.PatchUserRoles(ctx, models.User{Name: "pw-preserve", Roles: []string{}})

	ok, err = auth.VerifyUserPassword(ctx, "pw-preserve", "initial-pw")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("PatchPreservesPassword: password should still be valid")
	}
}
