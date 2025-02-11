package auth

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/auth"
)

func RequireScopeApiDecorator[Req any, Resp any](
	authStorage *auth.Storage,
	scope models.Scope,
	next func(context.Context, Req, *Resp) error,
) func(context.Context, Req, *Resp) error {
	return func(ctx context.Context, req Req, resp *Resp) error {
		user := GetContextUser(ctx)
		authorized, err := authStorage.VerifyUserPermission(
			ctx,
			user.Name,
			scope,
		)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"Failed to verify user permission",
				slog.String("channel", "api"),
				slog.String("error", err.Error()),
			)
			return status.Wrap(err, status.PermissionDenied)
		}

		if !authorized {
			return status.PermissionDenied
		}

		slog.DebugContext(
			ctx,
			"Authorized user",
			slog.String("channel", "api"),
			slog.String("user", user.Name),
			slog.String("scope", string(scope)),
		)

		return next(ctx, req, resp)
	}
}

func RequireScopesApiDecorator[Req any, Resp any](
	authStorage *auth.Storage,
	scopes []models.Scope,
	next func(context.Context, Req, *Resp) error,
) func(context.Context, Req, *Resp) error {
	if len(scopes) == 0 {
		return next
	} else if len(scopes) == 1 {
		return RequireScopeApiDecorator(authStorage, scopes[0], next)
	} else {
		return RequireScopesApiDecorator(
			authStorage,
			scopes[1:],
			RequireScopeApiDecorator(authStorage, scopes[0], next),
		)
	}
}
