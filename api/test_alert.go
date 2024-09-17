package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/data/logstorage"
)

type TestAlertRequest struct {
	Alert  string            `path:"alert" minLength:"1"`
	Record map[string]string `json:"record"`
}

type TestAlertResponse struct {
	Success bool `json:"success"`
}

func TestAlertUsecase(
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
				req TestAlertRequest,
				resp *TestAlertResponse,
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

				logentry := logstorage.NewLogEntry(req.Record)
				err = webhook.Call(ctx, logentry)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to call alert webhook",
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

	u.SetName("test_alert")
	u.SetTitle("Test Alert")
	u.SetDescription("Test alert")
	u.SetTags("tests")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}
