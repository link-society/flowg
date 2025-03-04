package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type ListAlertsRequest struct{}
type ListAlertsResponse struct {
	Success bool     `json:"success"`
	Alerts  []string `json:"alerts"`
}

func (ctrl *controller) ListAlertsUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_ALERTS,
			func(
				ctx context.Context,
				req ListAlertsRequest,
				resp *ListAlertsResponse,
			) error {
				alerts, err := ctrl.deps.ConfigStorage.ListAlerts(ctx)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to list alerts",
						slog.String("error", err.Error()),
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
