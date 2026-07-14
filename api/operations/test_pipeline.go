package operations

import (
	"context"
	"fmt"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/api/schemas"

	applog "link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// TestPipelineDeps lists the dependencies of [NewTestPipelineUsecase].
type TestPipelineDeps struct {
	fx.In

	AuthStorage    storage.AuthStorage
	PipelineRunner pipelines.Runner
}

// NewTestPipelineUsecase runs sample records through a pipeline definition and
// returns an execution trace without persisting anything.
//
// It lets authors observe how records flow through nodes before saving a
// pipeline. Callers must have the send-logs permission. The run is marked sensitive
// so its records are excluded from FlowG's own logs.
func NewTestPipelineUsecase(deps TestPipelineDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req schemas.TestPipelineRequest,
				resp *schemas.TestPipelineResponse,
			) error {
				tracer := pipelines.NodeTracer{
					Flow: req.Flow,
				}
				ctx = pipelines.WithTracer(ctx, &tracer)
				applog.MarkSensitive(ctx)

				for _, recordData := range req.Records {
					record := models.NewLogRecord(recordData)
					err := deps.PipelineRunner.Run(
						ctx,
						req.Pipeline,
						pipelines.DIRECT_ENTRYPOINT,
						record,
					)

					if err != nil {
						errMsg := fmt.Sprint(err)
						resp.Error = &errMsg
					}

					logger.DebugContext(
						ctx,
						"Ran test for",
						slog.String("pipeline", req.Pipeline),
					)
				}

				resp.Success = true
				resp.Trace = tracer.Trace

				return nil
			},
		),
	)

	u.SetName("test_pipeline")
	u.SetTitle("Test the pipeline")
	u.SetDescription("Test running structured logs through a pipeline")
	u.SetTags("tests")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewTestPipelineUsecase,
		http.MethodPost,
		"/api/v1/test/pipeline/{pipeline}",
	)
}
