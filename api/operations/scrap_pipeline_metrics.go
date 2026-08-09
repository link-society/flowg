package operations

import (
	"context"
	"errors"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/openapi-go"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/api/schemas"

	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// ScrapPipelineMetricsDeps lists the dependencies of [NewScrapPipelineMetricsUsecase].
type ScrapPipelineMetricsDeps struct {
	fx.In

	AuthStorage    storage.AuthStorage
	PipelineRunner pipelines.Runner
}

// NewScrapPipelineMetricsUsecase returns the metrics of a single pipeline.
//
// Callers must have the read-pipelines permission. Requesting an unknown pipeline
// yields a not-found error.
func NewScrapPipelineMetricsUsecase(deps ScrapPipelineMetricsDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req schemas.ScrapPipelineMetricsRequest,
				resp *schemas.ScrapPipelineMetricsResponse,
			) error {
				err := deps.PipelineRunner.ScrapMetrics(ctx, req.Pipeline, resp.Writer)
				if err != nil {
					if errors.Is(err, &pipelines.PipelineNotFoundError{}) {
						return status.Wrap(err, status.NotFound)
					}

					logger.ErrorContext(
						ctx,
						"Failed to scrap pipeline metrics",
						slog.String("pipeline", req.Pipeline),
						slog.String("error", err.Error()),
					)
					return status.Wrap(err, status.Internal)
				}

				return nil
			},
		),
	)

	u.SetName("scrap_pipeline_metrics")
	u.SetTitle("Scrap Pipeline Metrics")
	u.SetDescription("Prometheus Exporter for the pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}

// annotateScrapPipelineMetrics documents the scrap pipeline metrics response
// as plain text
func annotateScrapPipelineMetrics(oc openapi.OperationContext) error {
	contentUnits := oc.Response()
	for i, cu := range contentUnits {
		if cu.HTTPStatus == 200 {
			cu.ContentType = "text/plain"
			cu.Description = "Prometheus Exporter"
			cu.Format = "text"
		}

		contentUnits[i] = cu
	}

	return nil
}
func init() {
	routing.RegisterOperation(
		NewScrapPipelineMetricsUsecase,
		http.MethodGet,
		"/api/v1/pipelines/{pipeline}/metrics",
		routing.Annotated(annotateScrapPipelineMetrics),
	)
}
