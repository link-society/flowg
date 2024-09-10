package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
)

type DeleteAlertRequest struct {
	Alert string `path:"alert" minLength:"1"`
}

type DeleteAlertResponse struct {
	Success bool `json:"success"`
}

func DeleteAlertUsecase(
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
				req DeleteAlertRequest,
				resp *DeleteAlertResponse,
			) error {
				err := alertSys.Delete(req.Alert)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to delete alert",
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

	u.SetName("delete_alert")
	u.SetTitle("Delete Alert")
	u.SetDescription("Delete alert")
	u.SetTags("alerts")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
