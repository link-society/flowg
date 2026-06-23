package auth_test

import (
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/api/auth"
)

func TestContextWithUserRoundTrip(t *testing.T) {
	// Contract: a user bound with ContextWithUser is retrievable, unchanged,
	// through GetContextUser.
	user := &models.User{Name: "alice", Roles: []string{"admin"}}

	ctx := auth.ContextWithUser(context.Background(), user)
	got := auth.GetContextUser(ctx)

	require.NotNil(t, got)
	assert.Same(t, user, got)
}

func TestContextWithUserDoesNotMutateParent(t *testing.T) {
	// Contract: ContextWithUser derives a new context and leaves the parent
	// free of any identity.
	parent := context.Background()
	user := &models.User{Name: "bob", Roles: nil}

	_ = auth.ContextWithUser(parent, user)

	assert.Nil(t, parent.Value(auth.CONTEXT_USER))
}

func TestGetContextUserPanicsWhenMissing(t *testing.T) {
	// Contract: GetContextUser treats a missing identity as a programming
	// error and panics rather than returning a nil user.
	assert.Panics(t, func() {
		auth.GetContextUser(context.Background())
	})
}
