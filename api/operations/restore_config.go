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

// RestoreConfigDeps lists the dependencies of [NewRestoreConfigUsecase].
type RestoreConfigDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
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
				req schemas.RestoreConfigRequest,
				resp *schemas.RestoreConfigResponse,
			) error {
				defer req.Backup.Close()

				err := deps.ConfigStorage.Load(ctx, req.Backup)
				if err != nil {
					if errors.Is(err, kv.ErrNotSupported) {
						return status.Wrap(err, status.Unimplemented)
					}

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

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Unimplemented, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewRestoreConfigUsecase,
		http.MethodPost,
		"/api/v1/restore/config",
	)
}
