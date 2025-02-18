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
	"link-society.com/flowg/internal/storage/config"
)

type RestoreConfigRequest struct {
	Backup multipart.File `formData:"backup"`
}

type RestoreConfigResponse struct {
	Success bool `json:"success"`
}

func RestoreConfigUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopesApiDecorator(
			authStorage,
			[]models.Scope{
				models.SCOPE_WRITE_PIPELINES,
				models.SCOPE_WRITE_TRANSFORMERS,
				models.SCOPE_WRITE_ALERTS,
			},
			func(
				ctx context.Context,
				req RestoreConfigRequest,
				resp *RestoreConfigResponse,
			) error {
				defer req.Backup.Close()

				err := configStorage.Restore(ctx, req.Backup)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to restore configuration database",
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

	u.SetName("restore_config")
	u.SetTitle("Restore Configuration")
	u.SetDescription("Upload a full snapshot of the configuration database.")
	u.SetTags("backup")

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Internal)

	return u
}
