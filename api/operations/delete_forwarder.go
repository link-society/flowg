package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/models"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// DeleteForwarderDeps lists the dependencies of [NewDeleteForwarderUsecase].
type DeleteForwarderDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// DeleteForwarderRequest identifies the forwarder to remove.
type DeleteForwarderRequest struct {
	// Forwarder is the name of the forwarder to delete.
	Forwarder string `path:"forwarder" minLength:"1"`
}

// DeleteForwarderResponse reports the outcome of the deletion.
type DeleteForwarderResponse struct {
	// Success reports whether the forwarder was removed.
	Success bool `json:"success"`
}

// NewDeleteForwarderUsecase removes a forwarder.
//
// Callers must have the write-forwarders permission. Deleting an absent
// forwarder is treated as a success.
func NewDeleteForwarderUsecase(deps DeleteForwarderDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_FORWARDERS,
			func(
				ctx context.Context,
				req DeleteForwarderRequest,
				resp *DeleteForwarderResponse,
			) error {
				err := deps.ConfigStorage.DeleteForwarder(ctx, req.Forwarder)
				if err != nil {
					logger.ErrorContext(
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

func init() {
	routing.RegisterOperation(
		NewDeleteForwarderUsecase,
		http.MethodDelete,
		"/api/v1/forwarders/{forwarder}",
	)
}
