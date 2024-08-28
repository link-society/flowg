package auth

import (
	"context"
)

type authContextKey string

const CONTEXT_USER authContextKey = "user"

func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, CONTEXT_USER, user)
}

func GetContextUser(ctx context.Context) *User {
	return ctx.Value(CONTEXT_USER).(*User)
}
