package auth

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase/status"
)

func RequireScopeApiDecorator[Req any, Resp any](
	db *Database,
	scope Scope,
	next func(context.Context, Req, *Resp) error,
) func(context.Context, Req, *Resp) error {
	userSys := NewUserSystem(db)

	return func(ctx context.Context, req Req, resp *Resp) error {
		user := GetContextUser(ctx)
		authorized, err := userSys.VerifyUserPermission(
			user.Name,
			scope,
		)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"Failed to verify user permission",
				"channel", "api",
				"error", err.Error(),
			)
			return status.Wrap(err, status.PermissionDenied)
		}

		if !authorized {
			return status.PermissionDenied
		}

		slog.DebugContext(
			ctx,
			"Authorized user",
			"channel", "api",
			"user", user.Name,
			"scope", scope,
		)

		return next(ctx, req, resp)
	}
}
