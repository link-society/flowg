package api

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"link-society.com/flowg/internal/app/logging"
	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/engines/pipelines"
)

type TestPipelineRequest struct {
	Pipeline string              `path:"pipeline" minLength:"1"`
	Flow     models.FlowGraphV2  `json:"flow" required:"true"`
	Records  []map[string]string `json:"records" required:"true"`
}
type TestPipelineResponse struct {
	Success bool                  `json:"success"`
	Trace   []pipelines.NodeTrace `json:"trace"`
	Error   *string               `json:"error,omitempty"`
}

func (ctrl *controller) TestPipelineUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req TestPipelineRequest,
				resp *TestPipelineResponse,
			) error {
				tracer := pipelines.NodeTracer{
					Flow: req.Flow,
				}
				ctx = pipelines.WithTracer(ctx, &tracer)
				logging.MarkSensitive(ctx)

				for _, recordData := range req.Records {
					record := models.NewLogRecord(recordData)
					err := ctrl.deps.PipelineRunner.Run(
						ctx,
						req.Pipeline,
						pipelines.DIRECT_ENTRYPOINT,
						record,
					)

					if err != nil {
						errMsg := fmt.Sprint(err)
						resp.Error = &errMsg
					}

					ctrl.logger.DebugContext(
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
