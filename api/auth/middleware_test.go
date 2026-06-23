package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"link-society.com/flowg/internal/models"
	authStorage "link-society.com/flowg/internal/storage/auth"

	"link-society.com/flowg/api/auth"
)

func newCapturingHandler(captured **models.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*captured = auth.GetContextUser(r.Context())
		w.WriteHeader(http.StatusOK)
	})
}

func TestApiMiddlewareRejectsMissingCredential(t *testing.T) {
	// Contract: a request without an Authorization header never reaches the
	// next handler and is answered with an Unauthenticated status.
	mockStorage := authStorage.NewMockStorage().(*authStorage.MockStorage)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	handler := auth.ApiMiddleware(mockStorage)(next)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)

	assert.False(t, nextCalled, "next must not be invoked without a credential")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestApiMiddlewareAuthenticatesPersonalAccessToken(t *testing.T) {
	// Contract: a valid personal access token resolves the user and binds it
	// to the context handed to the next handler.
	mockStorage := authStorage.NewMockStorage().(*authStorage.MockStorage)
	user := &models.User{Name: "alice", Roles: []string{"admin"}}
	mockStorage.On("VerifyToken", mock.Anything, "pat_secret").
		Return(user, nil)

	var captured *models.User
	handler := auth.ApiMiddleware(mockStorage)(newCapturingHandler(&captured))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer pat_secret")
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, captured)
	assert.Same(t, user, captured)
	mockStorage.AssertExpectations(t)
}

func TestApiMiddlewareRejectsUnknownPersonalAccessToken(t *testing.T) {
	// Contract: a personal access token that resolves to no user is treated
	// as invalid and stops the chain.
	mockStorage := authStorage.NewMockStorage().(*authStorage.MockStorage)
	mockStorage.On("VerifyToken", mock.Anything, "pat_unknown").
		Return((*models.User)(nil), nil)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})
	handler := auth.ApiMiddleware(mockStorage)(next)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer pat_unknown")
	handler.ServeHTTP(rec, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestApiMiddlewareAuthenticatesJWT(t *testing.T) {
	// Contract: a valid JWT resolves the subject through the user store and
	// binds the user to the context handed to the next handler.
	mockStorage := authStorage.NewMockStorage().(*authStorage.MockStorage)
	user := &models.User{Name: "bob", Roles: []string{"reader"}}
	mockStorage.On("FetchUser", mock.Anything, "bob").
		Return(user, nil)

	token, err := auth.NewJWT("bob")
	require.NoError(t, err)

	var captured *models.User
	handler := auth.ApiMiddleware(mockStorage)(newCapturingHandler(&captured))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, captured)
	assert.Same(t, user, captured)
	mockStorage.AssertExpectations(t)
}

func TestApiMiddlewareRejectsUnknownScheme(t *testing.T) {
	// Contract: a credential that matches no supported scheme is rejected.
	mockStorage := authStorage.NewMockStorage().(*authStorage.MockStorage)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})
	handler := auth.ApiMiddleware(mockStorage)(next)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	handler.ServeHTTP(rec, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
