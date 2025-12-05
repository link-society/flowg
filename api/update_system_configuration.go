package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type UpdateSystemConfigurationRequest = models.SystemConfiguration

type UpdateSystemConfigurationResponse = struct {
	Success bool `json:"success"`
}

func (ctrl *controller) UpdateSystemConfigurationUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_SYSTEM_CONFIGURATION,
			func(
				ctx context.Context,
				req UpdateSystemConfigurationRequest,
				resp *UpdateSystemConfigurationResponse,
			) error {
				err := ctrl.deps.ConfigStorage.WriteSystemConfig(ctx, &req)
				if err != nil {
					ctrl.logger.ErrorContext(
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
