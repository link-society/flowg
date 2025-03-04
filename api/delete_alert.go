package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type DeleteAlertRequest struct {
	Alert string `path:"alert" minLength:"1"`
}

type DeleteAlertResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) DeleteAlertUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_ALERTS,
			func(
				ctx context.Context,
				req DeleteAlertRequest,
				resp *DeleteAlertResponse,
			) error {
				err := ctrl.deps.ConfigStorage.DeleteAlert(ctx, req.Alert)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to delete alert",
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

	u.SetName("delete_alert")
	u.SetTitle("Delete Alert")
	u.SetDescription("Delete alert")
	u.SetTags("alerts")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
