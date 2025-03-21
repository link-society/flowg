package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type DeleteForwarderRequest struct {
	Forwarder string `path:"forwarder" minLength:"1"`
}

type DeleteForwarderResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) DeleteForwarderUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_FORWARDERS,
			func(
				ctx context.Context,
				req DeleteForwarderRequest,
				resp *DeleteForwarderResponse,
			) error {
				err := ctrl.deps.ConfigStorage.DeleteForwarder(ctx, req.Forwarder)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to delete forwarder",
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

	u.SetName("delete_forwarder")
	u.SetTitle("Delete Forwarder")
	u.SetDescription("Delete forwarder")
	u.SetTags("forwarders")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
