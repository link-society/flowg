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

// BackupLogsDeps lists the dependencies of [NewBackupLogsUsecase].
type BackupLogsDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// BackupLogsRequest is empty: the whole log database is exported.
type BackupLogsRequest struct{}

// BackupLogsResponse streams the log database snapshot to the client.
//
// It embeds the writer so the snapshot is streamed as a file download rather
// than buffered in memory.
type BackupLogsResponse struct {
	usecase.OutputWithEmbeddedWriter
}

// NewBackupLogsUsecase streams a full snapshot of the log database as a
// downloadable file.
//
// It is the export half of the log backup story; the snapshot can later be
// reloaded with [NewRestoreLogsUsecase]. Callers must have the
// read-streams permission.
func NewBackupLogsUsecase(deps BackupLogsDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req BackupLogsRequest,
				resp *BackupLogsResponse,
			) error {
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Type", "application/octet-stream")
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Disposition", "attachment; filename=logs.db")
				resp.Writer.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")

				_, err := deps.LogStorage.Dump(ctx, resp.Writer, 0)
				resp.Writer.(http.Flusher).Flush()
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to backup logs database",
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

// annotateBackupLogs documents the backup response as a binary file download.
func annotateBackupLogs(oc openapi.OperationContext) error {
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
		NewBackupLogsUsecase,
		http.MethodGet,
		"/api/v1/backup/logs",
		routing.Annotated(annotateBackupLogs),
	)
}
