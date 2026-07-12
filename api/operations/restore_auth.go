package operations

import (
	"context"
	"errors"
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

	"link-society.com/flowg/internal/storage/generic/kv"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// RestoreAuthDeps lists the dependencies of [NewRestoreAuthUsecase].
type RestoreAuthDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// RestoreAuthRequest carries the authentication database snapshot to load.
type RestoreAuthRequest struct {
	// Backup is the uploaded snapshot, as produced by
	// [NewBackupAuthUsecase].
	Backup multipart.File `formData:"backup"`
}

// RestoreAuthResponse reports the outcome of the restore.
type RestoreAuthResponse struct {
	// Success reports whether the snapshot was loaded.
	Success bool `json:"success"`
}

// NewRestoreAuthUsecase loads a previously exported authentication database
// snapshot, replacing the current contents.
//
// It is the import counterpart to [NewBackupAuthUsecase]. Callers must
// have the write-ACLs permission.
func NewRestoreAuthUsecase(deps RestoreAuthDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req RestoreAuthRequest,
				resp *RestoreAuthResponse,
			) error {
				defer req.Backup.Close()

				err := deps.AuthStorage.Load(ctx, req.Backup)
				if err != nil {
					if errors.Is(err, kv.ErrNotSupported) {
						return status.Wrap(err, status.Unimplemented)
					}

					logger.ErrorContext(
						ctx,
						"Failed to restore authentication database",
						slog.String("error", err.Error()),
					)

					return status.Wrap(err, status.Internal)
				}

				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("restore_auth")
	u.SetTitle("Restore Authentication Database")
	u.SetDescription("Upload a full snapshot of the authentication database.")
	u.SetTags("backup")

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Unimplemented, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewRestoreAuthUsecase,
		http.MethodPost,
		"/api/v1/restore/auth",
	)
}
