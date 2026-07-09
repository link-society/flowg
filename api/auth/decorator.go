package auth

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/models"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// RequireScopeApiDecorator guards a use-case interactor so that it only runs
// for callers who have been granted a given permission scope.
//
// It exists to keep authorization out of the business logic: handlers declare
// the scope they need and remain unaware of how permissions are resolved. The
// returned interactor relies on the authenticated user being present in the
// context (see [GetContextUser]) and behaves as follows:
//
//   - the permission lookup fails: the error is wrapped as
//     [status.PermissionDenied] and next is not invoked;
//   - the user lacks the scope: [status.PermissionDenied] is returned and next
//     is not invoked;
//   - the user holds the scope: next is invoked and its result is returned
//     unchanged.
func RequireScopeApiDecorator[Req any, Resp any](
	authStorage storage.AuthStorage,
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

// RequireScopesApiDecorator guards a use-case interactor so that it only runs
// for callers who have been granted every one of the given permission scopes.
//
// It generalizes [RequireScopeApiDecorator] to the common case where an
// endpoint requires several permissions at once. Authorization is conjunctive:
// next is invoked only when all scopes are satisfied, and the first missing or
// unresolvable scope short-circuits with [status.PermissionDenied]. An empty
// scope list imposes no requirement and yields next unchanged.
func RequireScopesApiDecorator[Req any, Resp any](
	authStorage storage.AuthStorage,
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
