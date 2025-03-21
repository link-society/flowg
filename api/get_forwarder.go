package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type GetForwarderRequest struct {
	Forwarder string `path:"forwarder" minLength:"1"`
}

type GetForwarderResponse struct {
	Success   bool                `json:"success"`
	Forwarder *models.ForwarderV2 `json:"forwarder"`
}

func (ctrl *controller) GetForwarderUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_FORWARDERS,
			func(
				ctx context.Context,
				req GetForwarderRequest,
				resp *GetForwarderResponse,
			) error {
				forwarder, err := ctrl.deps.ConfigStorage.ReadForwarder(ctx, req.Forwarder)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to get forwarder",
						slog.String("forwarder", req.Forwarder),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				resp.Success = true
				resp.Forwarder = forwarder

				return nil
			},
		),
	)

	u.SetName("get_forwarder")
	u.SetTitle("Get Forwarder")
	u.SetDescription("Get forwarder")
	u.SetTags("forwarders")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound)

	return u
}
