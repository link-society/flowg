package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type TestAlertRequest struct {
	Alert  string            `path:"alert" minLength:"1"`
	Record map[string]string `json:"record"`
}

type TestAlertResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) TestAlertUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_ALERTS,
			func(
				ctx context.Context,
				req TestAlertRequest,
				resp *TestAlertResponse,
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

				logRecord := models.NewLogRecord(req.Record)
				err = webhook.Call(ctx, logRecord)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to call alert webhook",
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

	u.SetName("test_alert")
	u.SetTitle("Test Alert")
	u.SetDescription("Test alert")
	u.SetTags("tests")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}
