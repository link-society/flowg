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

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// RestoreConfigDeps lists the dependencies of [NewRestoreConfigUsecase].
type RestoreConfigDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// RestoreConfigRequest carries the configuration database snapshot to load.
type RestoreConfigRequest struct {
	// Backup is the uploaded snapshot, as produced by
	// [NewBackupConfigUsecase].
	Backup multipart.File `formData:"backup"`
}

// RestoreConfigResponse reports the outcome of the restore.
type RestoreConfigResponse struct {
	// Success reports whether the snapshot was loaded.
	Success bool `json:"success"`
}

// NewRestoreConfigUsecase loads a previously exported configuration database
// snapshot, replacing the current pipelines, transformers and forwarders.
//
// It is the import counterpart to [NewBackupConfigUsecase]. Callers
// must have write access to pipelines, transformers and forwarders.
func NewRestoreConfigUsecase(deps RestoreConfigDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopesApiDecorator(
			deps.AuthStorage,
			[]models.Scope{
				models.SCOPE_WRITE_PIPELINES,
				models.SCOPE_WRITE_TRANSFORMERS,
				models.SCOPE_WRITE_FORWARDERS,
			},
			func(
				ctx context.Context,
				req RestoreConfigRequest,
				resp *RestoreConfigResponse,
			) error {
				defer req.Backup.Close()

				err := deps.ConfigStorage.Load(ctx, req.Backup)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to restore configuration database",
						slog.String("error", err.Error()),
					)

					return status.Wrap(err, status.Internal)
				}

				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("restore_config")
	u.SetTitle("Restore Configuration")
	u.SetDescription("Upload a full snapshot of the configuration database.")
	u.SetTags("backup")

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewRestoreConfigUsecase,
		http.MethodPost,
		"/api/v1/restore/config",
	)
}
