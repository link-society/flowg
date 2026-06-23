package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/models"

	authStorage "link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
)

// UpdateSystemConfigurationDeps lists the dependencies of [NewUpdateSystemConfigurationUsecase].
type UpdateSystemConfigurationDeps struct {
	fx.In

	AuthStorage   authStorage.Storage
	ConfigStorage config.Storage
}

// UpdateSystemConfigurationRequest is the new system configuration to store. It
// aliases [models.SystemConfiguration] so the whole configuration is replaced
// in one request.
type UpdateSystemConfigurationRequest = models.SystemConfiguration

// UpdateSystemConfigurationResponse reports the outcome of the update.
type UpdateSystemConfigurationResponse = struct {
	// Success reports whether the configuration was persisted.
	Success bool `json:"success"`
}

// NewUpdateSystemConfigurationUsecase replaces the global system configuration.
//
// Callers must have the write-system-configuration permission.
func NewUpdateSystemConfigurationUsecase(deps UpdateSystemConfigurationDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_SYSTEM_CONFIGURATION,
			func(
				ctx context.Context,
				req UpdateSystemConfigurationRequest,
				resp *UpdateSystemConfigurationResponse,
			) error {
				err := deps.ConfigStorage.WriteSystemConfig(ctx, &req)
				if err != nil {
					logger.ErrorContext(
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
