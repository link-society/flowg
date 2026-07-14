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

	"link-society.com/flowg/internal/engines/forwarders"
	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// TestForwarderDeps lists the dependencies of [NewTestForwarderUsecase].
type TestForwarderDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// TestForwarderRequest carries the forwarder name and a sample record to send
// through it.
type TestForwarderRequest struct {
	// Forwarder is the name of the stored forwarder to exercise.
	Forwarder string `path:"forwarder" minLength:"1"`
	// Record is the sample log record to forward.
	Record map[string]string `json:"record" required:"true"`
}

// TestForwarderResponse reports whether the trial delivery succeeded.
type TestForwarderResponse struct {
	// Success reports whether the record was delivered.
	Success bool `json:"success"`
}

// NewTestForwarderUsecase delivers a sample record through a stored forwarder to
// verify it is reachable and correctly configured.
//
// It actually contacts the forwarder's destination, so a successful response
// confirms end-to-end connectivity. Callers must have the write-forwarders
// permission.
func NewTestForwarderUsecase(deps TestForwarderDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_FORWARDERS,
			func(
				ctx context.Context,
				req TestForwarderRequest,
				resp *TestForwarderResponse,
			) error {
				forwarder, err := deps.ConfigStorage.ReadForwarder(ctx, req.Forwarder)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to get forwarder",
						slog.String("forwarder", req.Forwarder),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				runtime, err := forwarders.NewRuntime(forwarder)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to create forwarder runtime",
						slog.String("forwarder", req.Forwarder),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				if err := runtime.Init(ctx); err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to initialize forwarder",
						slog.String("forwarder", req.Forwarder),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				defer func() {
					if err := runtime.Close(ctx); err != nil {
						logger.WarnContext(
							ctx,
							"Failed to shutdown forwarder",
							slog.String("forwarder", req.Forwarder),
							slog.String("error", err.Error()),
						)
					}
				}()

				logRecord := models.NewLogRecord(req.Record)
				if err := runtime.Call(ctx, logRecord); err != nil {
					logger.ErrorContext(
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

func init() {
	routing.RegisterOperation(
		NewTestForwarderUsecase,
		http.MethodPost,
		"/api/v1/test/forwarders/{forwarder}",
	)
}
