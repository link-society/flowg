package auth

import (
	"context"

	"link-society.com/flowg/internal/models"
)

type authContextKey string

// CONTEXT_USER is the context key under which the authenticated user is carried
// throughout the lifetime of a request.
const CONTEXT_USER authContextKey = "user"

// ContextWithUser binds an authenticated user to a request's context so that
// downstream handlers and decorators can make authorization decisions without
// re-authenticating.
//
// It is the single entry point for establishing the request identity: callers
// performing authentication (such as [ApiMiddleware]) use it, and consumers
// retrieve the value through [GetContextUser].
func ContextWithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, CONTEXT_USER, user)
}

// GetContextUser returns the authenticated user previously bound to ctx with
// [ContextWithUser].
//
// It assumes authentication has already happened upstream and therefore treats
// the absence of a user as a programming error: it panics rather than returning
// a nil identity that could silently bypass authorization checks. Only call it
// on a context that has passed through the authentication middleware.
func GetContextUser(ctx context.Context) *models.User {
	return ctx.Value(CONTEXT_USER).(*models.User)
}
