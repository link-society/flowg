package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type SaveAlertRequest struct {
	Alert   string           `path:"alert" minLength:"1"`
	Webhook models.WebhookV1 `json:"webhook"`
}

type SaveAlertResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) SaveAlertUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_ALERTS,
			func(
				ctx context.Context,
				req SaveAlertRequest,
				resp *SaveAlertResponse,
			) error {
				err := ctrl.deps.ConfigStorage.WriteAlert(ctx, req.Alert, &req.Webhook)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to save alert",
						slog.String("alert", req.Alert),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("save_alert")
	u.SetTitle("Save Alert")
	u.SetDescription("Save alert")
	u.SetTags("alerts")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
