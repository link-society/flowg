package api

import (
	"context"
	"log/slog"

	"net/http"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/models"
	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/log"
)

type BackupLogsRequest struct{}

type BackupLogsResponse struct {
	usecase.OutputWithEmbeddedWriter
}

func BackupLogsUsecase(
	authStorage *auth.Storage,
	logStorage *log.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req BackupLogsRequest,
				resp *BackupLogsResponse,
			) error {
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Type", "application/octet-stream")
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Disposition", "attachment; filename=logs.db")
				resp.Writer.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")

				err := logStorage.Backup(ctx, resp.Writer)
				resp.Writer.(http.Flusher).Flush()
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to backup logs database",
						slog.String("channel", "api"),
						slog.String("error", err.Error()),
					)

					return status.Wrap(err, status.Internal)
				}

				return nil
			},
		),
	)

	u.SetName("backup_logs")
	u.SetTitle("Backup Logs Database")
	u.SetDescription("Download a full snapshot of the logs database.")
	u.SetTags("backup")

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Internal)

	return u
}
