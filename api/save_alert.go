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

type SaveAlertRequest struct {
	Alert   string           `path:"alert" minLength:"1"`
	Webhook models.WebhookV1 `json:"webhook"`
}

type SaveAlertResponse struct {
	Success bool `json:"success"`
}

func SaveAlertUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_ALERTS,
			func(
				ctx context.Context,
				req SaveAlertRequest,
				resp *SaveAlertResponse,
			) error {
				err := configStorage.WriteAlert(ctx, req.Alert, &req.Webhook)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to save alert",
						slog.String("channel", "api"),
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
