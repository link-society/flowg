package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
)

type ListAlertsRequest struct{}
type ListAlertsResponse struct {
	Success bool     `json:"success"`
	Alerts  []string `json:"alerts"`
}

func ListAlertsUsecase(
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
				req ListAlertsRequest,
				resp *ListAlertsResponse,
			) error {
				alerts, err := alertSys.List()
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list alerts",
						"channel", "api",
						"error", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Alerts = alerts

				return nil
			},
		),
	)

	u.SetName("list_alerts")
	u.SetTitle("List Alerts")
	u.SetDescription("List alerts")
	u.SetTags("alerts")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
