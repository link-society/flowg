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

type DeleteAlertRequest struct {
	Alert string `path:"alert" minLength:"1"`
}

type DeleteAlertResponse struct {
	Success bool `json:"success"`
}

func DeleteAlertUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_ALERTS,
			func(
				ctx context.Context,
				req DeleteAlertRequest,
				resp *DeleteAlertResponse,
			) error {
				err := configStorage.DeleteAlert(ctx, req.Alert)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to delete alert",
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

	u.SetName("delete_alert")
	u.SetTitle("Delete Alert")
	u.SetDescription("Delete alert")
	u.SetTags("alerts")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
