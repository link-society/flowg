package auth

import (
	"context"

	"link-society.com/flowg/internal/models"
)

type authContextKey string

const CONTEXT_USER authContextKey = "user"

func ContextWithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, CONTEXT_USER, user)
}

func GetContextUser(ctx context.Context) *models.User {
	return ctx.Value(CONTEXT_USER).(*models.User)
}
