package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/swaggest/rest/request"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"
	"link-society.com/flowg/internal/utils/api/otlp"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/engines/pipelines"
)

type IngestOTLPRequest struct {
	Pipeline   string `path:"pipeline" minLength:"1"`
	logRecords []*models.LogRecord
}

func (ior *IngestOTLPRequest) LoadFromHTTPRequest(r *http.Request) (err error) {
	ior.logRecords, err = otlp.UnmarshalLogRecords(r)

	return err
}

var _ request.Loader = (*IngestOTLPRequest)(nil)

type IngestOTLPResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) IngestOTLPUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req *IngestOTLPRequest,
				resp *IngestOTLPResponse,
			) error {
				for _, logRecord := range req.logRecords {
					err := ctrl.deps.PipelineRunner.Run(
						ctx,
						req.Pipeline,
						pipelines.DIRECT_ENTRYPOINT,
						logRecord,
					)
					if err != nil {
						ctrl.logger.ErrorContext(
							ctx,
							"Failed to process log entry",
							slog.String("pipeline", req.Pipeline),
							slog.String("error", err.Error()),
						)

						resp.Success = false
						return status.Wrap(err, status.Internal)
					}
				}

				ctrl.logger.InfoContext(
					ctx,
					"Log entry processed",
					slog.String("pipeline", req.Pipeline),
				)
				resp.Success = true

				return nil
			},
		),
	)

	u.SetName(fmt.Sprintf("ingest_otlp logs"))
	u.SetTitle(fmt.Sprintf("Ingest OTLP logs"))

	u.SetDescription(fmt.Sprintf("Run otlp logs records through a pipeline"))
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
