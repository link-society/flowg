package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type SaveForwarderRequest struct {
	Forwarder string             `path:"forwarder" minLength:"1"`
	Config    models.ForwarderV2 `json:"config"`
}

type SaveForwarderResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) SaveForwarderUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_FORWARDERS,
			func(
				ctx context.Context,
				req SaveForwarderRequest,
				resp *SaveForwarderResponse,
			) error {
				err := ctrl.deps.ConfigStorage.WriteForwarder(ctx, req.Forwarder, &req.Config)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to save forwarder",
						slog.String("forwarder", req.Forwarder),
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

	u.SetName("save_forwarder")
	u.SetTitle("Save Forwarder")
	u.SetDescription("Save forwarder")
	u.SetTags("forwarders")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
