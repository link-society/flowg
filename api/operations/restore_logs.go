package operations

import (
	"context"
	"errors"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/api/schemas"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// RestoreLogsDeps lists the dependencies of [NewRestoreLogsUsecase].
type RestoreLogsDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
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
				req schemas.RestoreLogsRequest,
				resp *schemas.RestoreLogsResponse,
			) error {
				defer req.Backup.Close()

				err := deps.LogStorage.Load(ctx, req.Backup)
				if err != nil {
					if errors.Is(err, kv.ErrNotSupported) {
						return status.Wrap(err, status.Unimplemented)
					}

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

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Unimplemented, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewRestoreLogsUsecase,
		http.MethodPost,
		"/api/v1/restore/logs",
	)
}
