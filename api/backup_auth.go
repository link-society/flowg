package api

import (
	"context"
	"log/slog"

	"net/http"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/models"
	apiUtils "link-society.com/flowg/internal/utils/api"
)

type BackupAuthRequest struct{}

type BackupAuthResponse struct {
	usecase.OutputWithEmbeddedWriter
}

func (ctrl *controller) BackupAuthUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req BackupAuthRequest,
				resp *BackupAuthResponse,
			) error {
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Type", "application/octet-stream")
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Disposition", "attachment; filename=auth.db")
				resp.Writer.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")

				err := ctrl.deps.AuthStorage.Dump(ctx, resp.Writer, 0)
				resp.Writer.(http.Flusher).Flush()
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to backup authentication database",
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
