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

type RestoreLogsRequest struct {
	Backup multipart.File `formData:"backup"`
}

type RestoreLogsResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) RestoreLogsUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_STREAMS,
			func(
				ctx context.Context,
				req RestoreLogsRequest,
				resp *RestoreLogsResponse,
			) error {
				defer req.Backup.Close()

				err := ctrl.deps.LogStorage.Restore(ctx, req.Backup)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to restore logs database",
						slog.String("error", err.Error()),
					)

					return status.Wrap(err, status.Internal)
				}

				return nil
			},
		),
	)

	u.SetName("restore_logs")
	u.SetTitle("Restore Logs Database")
	u.SetDescription("Upload a full snapshot of the logs database.")
	u.SetTags("backup")

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Internal)

	return u
}
