package operations

import (
	"context"
	"log/slog"

	"mime/multipart"
	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage"
)

// RestoreLogsDeps lists the dependencies of [NewRestoreLogsUsecase].
type RestoreLogsDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// RestoreLogsRequest carries the log database snapshot to load.
type RestoreLogsRequest struct {
	// Backup is the uploaded snapshot, as produced by
	// [NewBackupLogsUsecase].
	Backup multipart.File `formData:"backup"`
}

// RestoreLogsResponse reports the outcome of the restore.
type RestoreLogsResponse struct {
	// Success reports whether the snapshot was loaded.
	Success bool `json:"success"`
}

// NewRestoreLogsUsecase loads a previously exported log database snapshot,
// replacing the current contents.
//
// It is the import counterpart to [NewBackupLogsUsecase]. Callers must
// have the write-streams permission.
func NewRestoreLogsUsecase(deps RestoreLogsDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_STREAMS,
			func(
				ctx context.Context,
				req RestoreLogsRequest,
				resp *RestoreLogsResponse,
			) error {
				defer req.Backup.Close()

				err := deps.LogStorage.Load(ctx, req.Backup)
				if err != nil {
					logger.ErrorContext(
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

func init() {
	routing.RegisterOperation(
		NewRestoreLogsUsecase,
		http.MethodPost,
		"/api/v1/restore/logs",
	)
}
