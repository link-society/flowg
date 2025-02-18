package api

import (
	"context"
	"log/slog"

	"mime/multipart"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/models"
	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/storage/auth"
)

type RestoreAuthRequest struct {
	Backup multipart.File `formData:"backup"`
}

type RestoreAuthResponse struct {
	Success bool `json:"success"`
}

func RestoreAuthUsecase(authStorage *auth.Storage) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req RestoreAuthRequest,
				resp *RestoreAuthResponse,
			) error {
				defer req.Backup.Close()

				err := authStorage.Restore(ctx, req.Backup)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to restore authentication database",
						slog.String("channel", "api"),
						slog.String("error", err.Error()),
					)

					return status.Wrap(err, status.Internal)
				}

				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("restore_auth")
	u.SetTitle("Restore Authentication Database")
	u.SetDescription("Upload a full snapshot of the authentication database.")
	u.SetTags("backup")

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Internal)

	return u
}
