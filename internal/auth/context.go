package auth

import (
	"context"
	"time"
)

type authContext struct {
	parent   context.Context
	username string
}

type authContextKey string

const CONTEXT_USERNAME authContextKey = "auth_context_username"

func ContextWithUsername(ctx context.Context, username string) context.Context {
	return &authContext{
		parent:   ctx,
		username: username,
	}
}

func (ac *authContext) Deadline() (deadline time.Time, ok bool) {
	return ac.parent.Deadline()
}

func (ac *authContext) Done() <-chan struct{} {
	return ac.parent.Done()
}

func (ac *authContext) Err() error {
	return ac.parent.Err()
}

func (ac *authContext) Value(key interface{}) interface{} {
	if key == CONTEXT_USERNAME && ac.username != "" {
		return ac.username
	}

	return ac.parent.Value(key)
}

func GetContextUser(ctx context.Context) string {
	res := ctx.Value(CONTEXT_USERNAME)
	if res == nil {
		return ""
	}

	return res.(string)
}
