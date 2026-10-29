package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"go.link-society.com/flowg/api/auth"
	"go.link-society.com/flowg/api/routing"
	"go.link-society.com/flowg/api/schemas"

	"go.link-society.com/flowg/internal/models"
	storage "go.link-society.com/flowg/internal/storage/interfaces"
)

// UpdateSystemConfigurationDeps lists the dependencies of [NewUpdateSystemConfigurationUsecase].
type UpdateSystemConfigurationDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// NewUpdateSystemConfigurationUsecase replaces the global system configuration.
//
// Callers must have the write-system-configuration permission.
func NewUpdateSystemConfigurationUsecase(deps UpdateSystemConfigurationDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.Logger,
			deps.AuthStorage,
			models.SCOPE_WRITE_SYSTEM_CONFIGURATION,
			func(
				ctx context.Context,
				req schemas.UpdateSystemConfigurationRequest,
				resp *schemas.UpdateSystemConfigurationResponse,
			) error {
				err := deps.ConfigStorage.WriteSystemConfig(ctx, &req)
				if err != nil {
					deps.Logger.ErrorContext(
						ctx,
						"Failed to write system configuration",
						slog.String("error", err.Error()),
					)

					return status.Wrap(err, status.Internal)
				}

				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("update_system_configuration")
	u.SetTitle("Update System configuration")
	u.SetDescription("Update System configuration")
	u.SetTags("config")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewUpdateSystemConfigurationUsecase,
		http.MethodPut,
		"/api/v1/system-configuration",
	)
}
