package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type TestForwarderRequest struct {
	Forwarder string            `path:"forwarder" minLength:"1"`
	Record    map[string]string `json:"record"`
}

type TestForwarderResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) TestForwarderUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_FORWARDERS,
			func(
				ctx context.Context,
				req TestForwarderRequest,
				resp *TestForwarderResponse,
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

				logRecord := models.NewLogRecord(req.Record)
				err = forwarder.Call(ctx, logRecord)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to call forwarder",
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

	u.SetName("test_forwarder")
	u.SetTitle("Test Forwarder")
	u.SetDescription("Test forwarder")
	u.SetTags("tests")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}
