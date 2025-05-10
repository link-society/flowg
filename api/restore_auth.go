package api

import (
	"context"
	"log/slog"

	"mime/multipart"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/models"
	apiUtils "link-society.com/flowg/internal/utils/api"
)

type RestoreAuthRequest struct {
	Backup multipart.File `formData:"backup"`
}

type RestoreAuthResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) RestoreAuthUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req RestoreAuthRequest,
				resp *RestoreAuthResponse,
			) error {
				defer req.Backup.Close()

				err := ctrl.deps.AuthStorage.Load(ctx, req.Backup)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to restore authentication database",
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
