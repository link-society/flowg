package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/logstorage"
	"link-society.com/flowg/internal/pipelines"
)

type IngestLogRequest struct {
	Pipeline string            `path:"pipeline" minLength:"1"`
	Record   map[string]string `json:"record"`
}
type IngestLogResponse struct {
	Success bool `json:"success"`
}

func IngestLogUsecase(pipelinesManager *pipelines.Manager) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req IngestLogRequest,
			resp *IngestLogResponse,
		) error {
			pipeline, err := pipelinesManager.GetPipeline(req.Pipeline)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to get pipeline",
					"channel", "api",
					"pipeline", req.Pipeline,
					"error", err.Error(),
				)
				return status.Wrap(err, status.NotFound)
			}

			entry := logstorage.NewLogEntry(req.Record)
			err = pipeline.Run(ctx, pipelinesManager, entry)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to process log entry",
					"channel", "api",
					"pipeline", req.Pipeline,
					"error", err.Error(),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			slog.InfoContext(
				ctx,
				"Log entry processed",
				"channel", "api",
				"pipeline", req.Pipeline,
			)
			resp.Success = true

			return nil
		},
	)

	u.SetName("ingest_log")
	u.SetTitle("Ingest Log")
	u.SetDescription("Run log record through a pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.NotFound, status.Internal)

	return u
}
