package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type GetAlertRequest struct {
	Alert string `path:"alert" minLength:"1"`
}

type GetAlertResponse struct {
	Success bool              `json:"success"`
	Webhook *models.WebhookV1 `json:"webhook"`
}

func (ctrl *controller) GetAlertUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_ALERTS,
			func(
				ctx context.Context,
				req GetAlertRequest,
				resp *GetAlertResponse,
			) error {
				webhook, err := ctrl.deps.ConfigStorage.ReadAlert(ctx, req.Alert)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to get alert",
						slog.String("alert", req.Alert),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				resp.Success = true
				resp.Webhook = webhook

				return nil
			},
		),
	)

	u.SetName("get_alert")
	u.SetTitle("Get Alert")
	u.SetDescription("Get alert")
	u.SetTags("alerts")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound)

	return u
}
