package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/alerting"
	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
)

type GetAlertRequest struct {
	Alert string `path:"alert" minLength:"1"`
}

type GetAlertResponse struct {
	Success bool              `json:"success"`
	Webhook *alerting.Webhook `json:"webhook"`
}

func GetAlertUsecase(
	authDb *auth.Database,
	configStorage *config.Storage,
) usecase.Interactor {
	alertSys := config.NewAlertSystem(configStorage)

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_ALERTS,
			func(
				ctx context.Context,
				req GetAlertRequest,
				resp *GetAlertResponse,
			) error {
				webhook, err := alertSys.Read(req.Alert)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get alert",
						"channel", "api",
						"alert", req.Alert,
						"error", err.Error(),
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
