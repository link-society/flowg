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

type TestAlertRequest struct {
	Alert  string            `path:"alert" minLength:"1"`
	Record map[string]string `json:"record"`
}

type TestAlertResponse struct {
	Success bool `json:"success"`
}

func TestAlertUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_ALERTS,
			func(
				ctx context.Context,
				req TestAlertRequest,
				resp *TestAlertResponse,
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

				logRecord := models.NewLogRecord(req.Record)
				err = webhook.Call(ctx, logRecord)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to call alert webhook",
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

	u.SetName("test_alert")
	u.SetTitle("Test Alert")
	u.SetDescription("Test alert")
	u.SetTags("tests")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}
