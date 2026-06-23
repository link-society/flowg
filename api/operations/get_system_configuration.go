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

// GetSystemConfigurationDeps lists the dependencies of [NewGetSystemConfigurationUsecase].
type GetSystemConfigurationDeps struct {
	fx.In

	AuthStorage   authStorage.Storage
	ConfigStorage config.Storage
}

// GetSystemConfigurationRequest is empty: the system configuration is global.
type GetSystemConfigurationRequest struct{}

// GetSystemConfigurationResponse carries the current system configuration.
type GetSystemConfigurationResponse = struct {
	// Success reports whether the configuration was returned.
	Success bool `json:"success"`
	// Configuration is the current global system configuration.
	Configuration models.SystemConfiguration `json:"configuration"`
}

// NewGetSystemConfigurationUsecase returns the global system configuration.
//
// Callers must have the read-system-configuration permission.
func NewGetSystemConfigurationUsecase(deps GetSystemConfigurationDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_SYSTEM_CONFIGURATION,
			func(
				ctx context.Context,
				req GetSystemConfigurationRequest,
				resp *GetSystemConfigurationResponse,
			) error {
				conf, err := deps.ConfigStorage.ReadSystemConfig(ctx)
				if err != nil {
					logger.ErrorContext(
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
