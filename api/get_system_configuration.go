package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type GetSystemConfigurationRequest struct{}

type GetSystemConfigurationResponse = struct {
	Success       bool                       `json:"success"`
	Configuration models.SystemConfiguration `json:"configuration"`
}

func (ctrl *controller) GetSystemConfigurationUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_SYSTEM_CONFIGURATION,
			func(
				ctx context.Context,
				req GetSystemConfigurationRequest,
				resp *GetSystemConfigurationResponse,
			) error {
				conf, err := ctrl.deps.ConfigStorage.ReadSystemConfig(ctx)
				if err != nil {
					ctrl.logger.ErrorContext(
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
