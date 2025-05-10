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

type RestoreConfigRequest struct {
	Backup multipart.File `formData:"backup"`
}

type RestoreConfigResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) RestoreConfigUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopesApiDecorator(
			ctrl.deps.AuthStorage,
			[]models.Scope{
				models.SCOPE_WRITE_PIPELINES,
				models.SCOPE_WRITE_TRANSFORMERS,
				models.SCOPE_WRITE_FORWARDERS,
			},
			func(
				ctx context.Context,
				req RestoreConfigRequest,
				resp *RestoreConfigResponse,
			) error {
				defer req.Backup.Close()

				err := ctrl.deps.ConfigStorage.Load(ctx, req.Backup)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to restore configuration database",
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
