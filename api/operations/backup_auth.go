package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/openapi-go"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/models"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// BackupAuthDeps lists the dependencies of [NewBackupAuthUsecase].
type BackupAuthDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// BackupAuthRequest is empty: the whole authentication database is exported.
type BackupAuthRequest struct{}

// BackupAuthResponse streams the authentication database snapshot to the client.
//
// It embeds the writer so the snapshot is streamed as a file download rather
// than buffered in memory.
type BackupAuthResponse struct {
	usecase.OutputWithEmbeddedWriter
}

// NewBackupAuthUsecase streams a full snapshot of the authentication database as a
// downloadable file.
//
// It is the export half of FlowG's backup story; the snapshot can later be
// reloaded with [NewRestoreAuthUsecase]. Callers must have the
// read-ACLs permission.
func NewBackupAuthUsecase(deps BackupAuthDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req BackupAuthRequest,
				resp *BackupAuthResponse,
			) error {
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Type", "application/octet-stream")
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Disposition", "attachment; filename=auth.db")
				resp.Writer.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")

				_, err := deps.AuthStorage.Dump(ctx, resp.Writer, 0)
				resp.Writer.(http.Flusher).Flush()
				if err != nil {
					logger.ErrorContext(
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

// annotateBackupAuth documents the backup response as a binary file download.
func annotateBackupAuth(oc openapi.OperationContext) error {
	contentUnits := oc.Response()
	for i, cu := range contentUnits {
		if cu.HTTPStatus == 200 {
			cu.ContentType = "application/octet-stream"
			cu.Description = "Binary file"
			cu.Format = "Binary file"
		}

		contentUnits[i] = cu
	}

	return nil
}

func init() {
	routing.RegisterOperation(
		NewBackupAuthUsecase,
		http.MethodGet,
		"/api/v1/backup/auth",
		routing.Annotated(annotateBackupAuth),
	)
}
