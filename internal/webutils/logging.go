package webutils

import (
	"context"
	"log/slog"

	"link-society.com/flowg/internal/data/auth"
)

func LogError(ctx context.Context, message string, err error, args ...any) {
	user := auth.GetContextUser(ctx)
	args = append(args,
		"channel", "web",
		"auth.logged_user", user.Name,
		"error", err.Error(),
	)
	slog.ErrorContext(ctx, message, args...)
}
