package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"google.golang.org/protobuf/proto"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"

	logsmodels "go.opentelemetry.io/proto/otlp/logs/v1"

	"link-society.com/flowg/internal/storage/config"
	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/engines/pipelines"
)

// OTLPDataToLogRecordFields convert a struct to a map[string]string
func OTLPDataToLogRecordFields(obj interface{}) (map[string]string, error) {
	out := make(map[string]string)
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %T", obj)
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name
		tag := fieldType.Tag.Get("json")
		if tag != "" && tag != "-" {
			if commaIndex := len(tag); commaIndex > 0 && tag[commaIndex-1] == ',' {
				tag = tag[:commaIndex-1]
			}
			fieldName = tag
		}

		b, err := json.Marshal(field.Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field %s: %w", fieldName, err)
		}
		out[fieldName] = string(b)
	}
	return out, nil
}

func LogToLogRecord(logRecord *logsmodels.LogRecord) (result *models.LogRecord, err error) {
	result = &models.LogRecord{
		Timestamp: time.Unix(0, int64(logRecord.TimeUnixNano)),
	}

	result.Fields, err = OTLPDataToLogRecordFields(logRecord)

	return
}

type ContentType string

const (
	ProtoContentType ContentType = "application/x-protobuf"
	JsonContentType  ContentType = "application/json"
)

func SendToPipeline(ctx context.Context, logRecords []*models.LogRecord, pipelineName string, pipelineRunner *pipelines.Runner, configStorage *config.Storage, logger *slog.Logger) error {
	for _, logRecord := range logRecords {

		err := pipelineRunner.Run(
			ctx,
			pipelineName,
			pipelines.SYSLOG_ENTRYPOINT,
			logRecord,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

type IngestOTLPRequest struct {
	Pipeline   string `path:"pipeline" minLength:"1"`
	logRecords []*models.LogRecord
}

func (ior *IngestOTLPRequest) LoadFromHTTPRequest(r *http.Request) (err error) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return err
	}

	contentType := ContentType(r.Header.Get("Content-Type"))
	message := collectlogs.ExportLogsServiceRequest{}

	switch contentType {
	case ProtoContentType:
		err := proto.Unmarshal(body, &message)
		if err != nil {
			return err
		}

	case JsonContentType:
		err := json.Unmarshal(body, &message)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	for _, resourceLogs := range message.GetResourceLogs() {
		for _, scopeLogs := range resourceLogs.GetScopeLogs() {
			for _, logRecord := range scopeLogs.GetLogRecords() {
				logRecordModel, err := LogToLogRecord(logRecord)
				if err != nil {
					return err
				}
				ior.logRecords = append(ior.logRecords, logRecordModel)
			}
		}
	}

	return nil
}

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
				req IngestOTLPRequest,
				resp *IngestOTLPResponse,
			) error {
				err := SendToPipeline(ctx, req.logRecords, req.Pipeline, ctrl.deps.PipelineRunner, ctrl.deps.ConfigStorage, ctrl.logger)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to send log records to pipelines",
						slog.String("error", err.Error()),
						slog.String("pipeline", req.Pipeline),
					)
					return status.Wrap(err, status.Internal)
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
