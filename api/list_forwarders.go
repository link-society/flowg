package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type ListForwardersRequest struct{}
type ListForwardersResponse struct {
	Success    bool     `json:"success"`
	Forwarders []string `json:"forwarders"`
}

func (ctrl *controller) ListForwardersUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_FORWARDERS,
			func(
				ctx context.Context,
				req ListForwardersRequest,
				resp *ListForwardersResponse,
			) error {
				forwarders, err := ctrl.deps.ConfigStorage.ListForwarders(ctx)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to list forwarders",
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Forwarders = forwarders

				return nil
			},
		),
	)

	u.SetName("list_forwarders")
	u.SetTitle("List Forwarders")
	u.SetDescription("List forwarders")
	u.SetTags("forwarders")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
