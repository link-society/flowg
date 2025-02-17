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
)

type BackupAuthRequest struct{}

type BackupAuthResponse struct {
	usecase.OutputWithEmbeddedWriter
}

func BackupAuthUsecase(authStorage *auth.Storage) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req BackupAuthRequest,
				resp *BackupAuthResponse,
			) error {
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Type", "application/octet-stream")
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Disposition", "attachment; filename=auth.db")
				resp.Writer.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")

				err := authStorage.Backup(ctx, resp.Writer)
				resp.Writer.(http.Flusher).Flush()
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to backup authentication database",
						slog.String("channel", "api"),
						slog.String("error", err.Error()),
					)

					return status.Wrap(err, status.Internal)
				}

				return nil
			},
		),
	)

	u.SetName("backup_auth")
	u.SetTitle("Backup Authentication Database")
	u.SetDescription("Download a full snapshot of the authentication database.")
	u.SetTags("backup")

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Internal)

	return u
}
