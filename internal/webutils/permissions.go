package webutils

import (
	"context"
	"log/slog"

	"link-society.com/flowg/internal/data/auth"
)

type permissionSystem struct {
	permissions auth.Permissions
}

type permissionContextKey string

const CONTEXT_PERMISSION_SYSTEM permissionContextKey = "permission_system"

func WithPermissionSystem(
	ctx context.Context,
	userSys *auth.UserSystem,
) context.Context {
	user := auth.GetContextUser(ctx)
	scopes, err := userSys.ListUserScopes(user.Name)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"Failed to fetch user permissions",
			"channel", "web",
			"user", user.Name,
			"error", err.Error(),
		)

		NotifyError(ctx, "Could not fetch user permissions")
		scopes = []auth.Scope{}
	}

	sys := &permissionSystem{
		permissions: auth.PermissionsFromScopes(scopes),
	}

	return context.WithValue(ctx, CONTEXT_PERMISSION_SYSTEM, sys)
}

func Permissions(ctx context.Context) auth.Permissions {
	sys := ctx.Value(CONTEXT_PERMISSION_SYSTEM).(*permissionSystem)
	return sys.permissions
}
