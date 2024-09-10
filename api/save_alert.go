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

type SaveAlertRequest struct {
	Alert   string           `path:"alert" minLength:"1"`
	Webhook alerting.Webhook `json:"webhook"`
}

type SaveAlertResponse struct {
	Success bool `json:"success"`
}

func SaveAlertUsecase(
	authDb *auth.Database,
	configStorage *config.Storage,
) usecase.Interactor {
	alertSys := config.NewAlertSystem(configStorage)

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_WRITE_ALERTS,
			func(
				ctx context.Context,
				req SaveAlertRequest,
				resp *SaveAlertResponse,
			) error {
				err := alertSys.Write(req.Alert, &req.Webhook)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to save alert",
						"channel", "api",
						"alert", req.Alert,
						"error", err.Error(),
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
