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
	"link-society.com/flowg/internal/storage/log"
)

type RestoreLogsRequest struct {
	Backup multipart.File `formData:"backup"`
}

type RestoreLogsResponse struct {
	Success bool `json:"success"`
}

func RestoreLogsUsecase(
	authStorage *auth.Storage,
	logStorage *log.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_STREAMS,
			func(
				ctx context.Context,
				req RestoreLogsRequest,
				resp *RestoreLogsResponse,
			) error {
				defer req.Backup.Close()

				err := logStorage.Restore(ctx, req.Backup)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to restore logs database",
						slog.String("channel", "api"),
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
