package auth_test

import (
	"context"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/cmd/flowg-server/logging"
	"link-society.com/flowg/internal/models"

	badgerauth "link-society.com/flowg/internal/storage/backends/badger/concrete/auth"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

func newAuthStorage(t *testing.T) (context.Context, storage.AuthStorage) {
	t.Helper()
	logging.Discard()

	opts := badgerauth.DefaultOptions()
	opts.InMemory = true

	var authStorage storage.AuthStorage

	app := fxtest.New(
		t,
		badgerauth.NewStorage(opts),
		fx.Populate(&authStorage),
		fx.NopLogger,
	)
	app.RequireStart()
	t.Cleanup(app.RequireStop)

	return t.Context(), authStorage
}

// TestTokenLifecycle exercises PAT creation, O(1) verification, revocation, and
// the reverse-index cleanup on both DeleteToken and DeleteUser.
func TestTokenLifecycle(t *testing.T) {
	ctx, authStorage := newAuthStorage(t)

	if err := authStorage.SaveUser(ctx, models.User{Name: "alice", Roles: []string{}}, "s3cret"); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	token, tokenUUID, err := authStorage.CreateToken(ctx, "alice")
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	// A valid token resolves to its owner.
	user, err := authStorage.VerifyToken(ctx, token)
	if err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}
	if user == nil || user.Name != "alice" {
		t.Fatalf("expected token to resolve to alice, got %+v", user)
	}

	// An unknown token resolves to nothing, without error.
	unknown, err := authStorage.VerifyToken(ctx, "pat_unknown")
	if err != nil {
		t.Fatalf("unexpected error verifying unknown token: %v", err)
	}
	if unknown != nil {
		t.Fatalf("expected unknown token to resolve to nil, got %+v", unknown)
	}

	// The token shows up in its owner's list.
	uuids, err := authStorage.ListTokens(ctx, "alice")
	if err != nil {
		t.Fatalf("failed to list tokens: %v", err)
	}
	if len(uuids) != 1 || uuids[0] != tokenUUID {
		t.Fatalf("expected [%s], got %v", tokenUUID, uuids)
	}

	// Revoking the token makes it (and its index entry) unusable.
	if err := authStorage.DeleteToken(ctx, "alice", tokenUUID); err != nil {
		t.Fatalf("failed to delete token: %v", err)
	}
	revoked, err := authStorage.VerifyToken(ctx, token)
	if err != nil {
		t.Fatalf("failed to verify revoked token: %v", err)
	}
	if revoked != nil {
		t.Fatalf("expected revoked token to resolve to nil, got %+v", revoked)
	}

	// Deleting the user must also drop the reverse index of their tokens, so a
	// still-held token can never resolve to the deleted account.
	token2, _, err := authStorage.CreateToken(ctx, "alice")
	if err != nil {
		t.Fatalf("failed to create second token: %v", err)
	}
	if err := authStorage.DeleteUser(ctx, "alice"); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}
	orphan, err := authStorage.VerifyToken(ctx, token2)
	if err != nil {
		t.Fatalf("failed to verify token of deleted user: %v", err)
	}
	if orphan != nil {
		t.Fatalf("expected token of deleted user to resolve to nil, got %+v", orphan)
	}
}
