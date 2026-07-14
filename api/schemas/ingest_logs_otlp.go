package schemas

import (
	"fmt"
	"log/slog"

	"compress/gzip"
	"io"
	"net/http"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/otlp"
)

// IngestLogsOTLPRequest carries an OpenTelemetry logs export to push through a
// pipeline.
//
// It implements [request.Loader] because the OTLP payload may be protobuf or
// JSON and optionally gzip-compressed, which the generic decoder cannot handle
// on its own.
type IngestLogsOTLPRequest struct {
	// Pipeline is the name of the pipeline to run the records through.
	Pipeline string `path:"pipeline" minLength:"1"`
	// ContentEncoding is the payload's transfer encoding; only gzip is accepted.
	ContentEncoding string `header:"Content-Encoding" enum:"gzip" required:"false"`
	// LogRecords holds the records decoded from the OTLP payload.
	LogRecords []*models.LogRecord

	collectlogs.ExportLogsServiceRequest
}

// LoadFromHTTPRequest decodes the OTLP payload from the raw HTTP request,
// handling optional gzip compression and both protobuf and JSON encodings.
//
// It populates the request's pipeline name and decoded records so the usecase
// can stay agnostic of the wire format.
func (ior *IngestLogsOTLPRequest) LoadFromHTTPRequest(r *http.Request) error {
	ior.Pipeline = r.PathValue("pipeline")
	if ior.Pipeline == "" {
		return fmt.Errorf("pipeline is required")
	}
	defer r.Body.Close()

	slog.InfoContext(
		r.Context(),
		"Parsing OpenTelemetry message",
		slog.String("otlp.content-type", r.Header.Get("Content-Type")),
		slog.String("otlp.content-encoding", r.Header.Get("Content-Encoding")),
	)

	var body []byte

	ior.ContentEncoding = r.Header.Get("Content-Encoding")
	switch ior.ContentEncoding {
	case "gzip":
		// decompress body
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gz.Close()

		data, err := io.ReadAll(gz)
		if err != nil {
			return fmt.Errorf("failed to read gzip body: %w", err)
		}

		body = data

	case "":
		data, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("failed to read raw body: %w", err)
		}

		body = data

	default:
		return fmt.Errorf("unsupported content encoding: %s", ior.ContentEncoding)
	}

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/x-protobuf":
		logRecords, err := otlp.UnmarshalProtobuf(body)
		if err != nil {
			return fmt.Errorf("failed to unmarshal protobuf: %w", err)
		}

		ior.LogRecords = logRecords

	case "application/json":
		logRecords, err := otlp.UnmarshalJSON(body)
		if err != nil {
			return fmt.Errorf("failed to unmarshal json: %w", err)
		}

		ior.LogRecords = logRecords

	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	return nil
}

// IngestLogsOTLPResponse reports how many records were processed.
type IngestLogsOTLPResponse struct {
	// Success reports whether every record was processed.
	Success bool `json:"success"`
	// ProcessedCount is the number of records that ran through the pipeline.
	ProcessedCount int `json:"processed_count"`
}
