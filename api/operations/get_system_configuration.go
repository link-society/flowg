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

	storage "go.link-society.com/flowg/internal/storage/interfaces"
	"go.link-society.com/flowg/pkg/models"
)

// GetSystemConfigurationDeps lists the dependencies of [NewGetSystemConfigurationUsecase].
type GetSystemConfigurationDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// NewGetSystemConfigurationUsecase returns the global system configuration.
//
// Callers must have the read-system-configuration permission.
func NewGetSystemConfigurationUsecase(deps GetSystemConfigurationDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.Logger,
			deps.AuthStorage,
			models.SCOPE_READ_SYSTEM_CONFIGURATION,
			func(
				ctx context.Context,
				req schemas.GetSystemConfigurationRequest,
				resp *schemas.GetSystemConfigurationResponse,
			) error {
				conf, err := deps.ConfigStorage.ReadSystemConfig(ctx)
				if err != nil {
					deps.Logger.ErrorContext(
						ctx,
						"Failed to read system configuration",
						slog.String("error", err.Error()),
					)

					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Configuration = *conf

				return nil
			},
		),
	)

	u.SetName("get_system_configuration")
	u.SetTitle("Get System configuration")
	u.SetDescription("Get System configuration")
	u.SetTags("config")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewGetSystemConfigurationUsecase,
		http.MethodGet,
		"/api/v1/system-configuration",
	)
}
