package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
)

type GetAlertRequest struct {
	Alert string `path:"alert" minLength:"1"`
}

type GetAlertResponse struct {
	Success bool            `json:"success"`
	Webhook *models.Webhook `json:"webhook"`
}

func GetAlertUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_ALERTS,
			func(
				ctx context.Context,
				req GetAlertRequest,
				resp *GetAlertResponse,
			) error {
				webhook, err := configStorage.ReadAlert(ctx, req.Alert)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get alert",
						slog.String("channel", "api"),
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
