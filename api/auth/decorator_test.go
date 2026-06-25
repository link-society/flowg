package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"

	"link-society.com/flowg/api/auth"
)

type req struct{}
type resp struct{}

func authorizedContext() context.Context {
	user := &models.User{Name: "alice", Roles: []string{"admin"}}
	return auth.ContextWithUser(context.Background(), user)
}

func TestRequireScopeAllowsAuthorizedCaller(t *testing.T) {
	// Contract: when the caller holds the scope, next is invoked and its
	// result is returned unchanged.
	mockStorage := storage.NewMockAuthStorage().(*storage.MockAuthStorage)
	mockStorage.On(
		"VerifyUserPermission",
		mock.Anything, "alice", models.SCOPE_READ_PIPELINES,
	).Return(true, nil)

	called := false
	sentinel := errors.New("from next")
	next := func(ctx context.Context, r req, w *resp) error {
		called = true
		return sentinel
	}

	decorated := auth.RequireScopeApiDecorator(
		mockStorage, models.SCOPE_READ_PIPELINES, next,
	)
	err := decorated(authorizedContext(), req{}, &resp{})

	assert.True(t, called, "next must be invoked for an authorized caller")
	assert.Same(t, sentinel, err, "next's result must be returned unchanged")
	mockStorage.AssertExpectations(t)
}

func TestRequireScopeRejectsUnauthorizedCaller(t *testing.T) {
	// Contract: when the caller lacks the scope, next is not invoked and an
	// error is returned.
	mockStorage := storage.NewMockAuthStorage().(*storage.MockAuthStorage)
	mockStorage.On(
		"VerifyUserPermission",
		mock.Anything, "alice", models.SCOPE_READ_PIPELINES,
	).Return(false, nil)

	called := false
	next := func(ctx context.Context, r req, w *resp) error {
		called = true
		return nil
	}

	decorated := auth.RequireScopeApiDecorator(
		mockStorage, models.SCOPE_READ_PIPELINES, next,
	)
	err := decorated(authorizedContext(), req{}, &resp{})

	assert.False(t, called, "next must not be invoked for an unauthorized caller")
	require.Error(t, err)
}

func TestRequireScopeForwardsLookupFailure(t *testing.T) {
	// Contract: when permission resolution fails, the error is surfaced and
	// next is not invoked.
	mockStorage := storage.NewMockAuthStorage().(*storage.MockAuthStorage)
	mockStorage.On(
		"VerifyUserPermission",
		mock.Anything, "alice", models.SCOPE_READ_PIPELINES,
	).Return(false, errors.New("storage down"))

	called := false
	next := func(ctx context.Context, r req, w *resp) error {
		called = true
		return nil
	}

	decorated := auth.RequireScopeApiDecorator(
		mockStorage, models.SCOPE_READ_PIPELINES, next,
	)
	err := decorated(authorizedContext(), req{}, &resp{})

	assert.False(t, called, "next must not be invoked when the lookup fails")
	require.Error(t, err)
}

func TestRequireScopesWithoutScopesRunsNextUnchanged(t *testing.T) {
	// Contract: an empty scope list imposes no requirement; next runs and no
	// permission lookups are performed.
	mockStorage := storage.NewMockAuthStorage().(*storage.MockAuthStorage)

	called := false
	next := func(ctx context.Context, r req, w *resp) error {
		called = true
		return nil
	}

	decorated := auth.RequireScopesApiDecorator(
		mockStorage, nil, next,
	)
	err := decorated(authorizedContext(), req{}, &resp{})

	assert.True(t, called, "next must be invoked when no scope is required")
	require.NoError(t, err)
	mockStorage.AssertNotCalled(t, "VerifyUserPermission", mock.Anything, mock.Anything, mock.Anything)
}

func TestRequireScopesRequiresEveryScope(t *testing.T) {
	// Contract: authorization is conjunctive; next runs only when all scopes
	// are satisfied.
	mockStorage := storage.NewMockAuthStorage().(*storage.MockAuthStorage)
	mockStorage.On(
		"VerifyUserPermission",
		mock.Anything, "alice", models.SCOPE_READ_PIPELINES,
	).Return(true, nil)
	mockStorage.On(
		"VerifyUserPermission",
		mock.Anything, "alice", models.SCOPE_WRITE_PIPELINES,
	).Return(true, nil)

	called := false
	next := func(ctx context.Context, r req, w *resp) error {
		called = true
		return nil
	}

	decorated := auth.RequireScopesApiDecorator(
		mockStorage,
		[]models.Scope{models.SCOPE_READ_PIPELINES, models.SCOPE_WRITE_PIPELINES},
		next,
	)
	err := decorated(authorizedContext(), req{}, &resp{})

	assert.True(t, called, "next must be invoked when all scopes are satisfied")
	require.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestRequireScopesRejectsWhenOneScopeMissing(t *testing.T) {
	// Contract: a single missing scope short-circuits and prevents next from
	// running.
	mockStorage := storage.NewMockAuthStorage().(*storage.MockAuthStorage)
	mockStorage.On(
		"VerifyUserPermission",
		mock.Anything, "alice", models.SCOPE_READ_PIPELINES,
	).Return(false, nil)
	// The write scope may or may not be evaluated depending on short-circuit
	// order; allow it without requiring it.
	mockStorage.On(
		"VerifyUserPermission",
		mock.Anything, "alice", models.SCOPE_WRITE_PIPELINES,
	).Return(true, nil).Maybe()

	called := false
	next := func(ctx context.Context, r req, w *resp) error {
		called = true
		return nil
	}

	decorated := auth.RequireScopesApiDecorator(
		mockStorage,
		[]models.Scope{models.SCOPE_READ_PIPELINES, models.SCOPE_WRITE_PIPELINES},
		next,
	)
	err := decorated(authorizedContext(), req{}, &resp{})

	assert.False(t, called, "next must not be invoked when a scope is missing")
	require.Error(t, err)
}
