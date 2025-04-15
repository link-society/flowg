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
	collectmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	collecttraces "go.opentelemetry.io/proto/otlp/collector/trace/v1"

	logsmodels "go.opentelemetry.io/proto/otlp/logs/v1"
	metricsmodels "go.opentelemetry.io/proto/otlp/metrics/v1"
	tracesmodels "go.opentelemetry.io/proto/otlp/trace/v1"

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

func MetricToLogRecord(metric *metricsmodels.Metric) (result *models.LogRecord, err error) {
	result = &models.LogRecord{
		Timestamp: time.Now(),
	}
	result.Fields, err = OTLPDataToLogRecordFields(metric)

	return
}

func LogToLogRecord(logRecord *logsmodels.LogRecord) (result *models.LogRecord, err error) {
	result = &models.LogRecord{
		Timestamp: time.Unix(0, int64(logRecord.TimeUnixNano)),
	}

	result.Fields, err = OTLPDataToLogRecordFields(logRecord)

	return
}

func SpanToLogRecord(span *tracesmodels.Span) (result *models.LogRecord, err error) {
	result = &models.LogRecord{
		Timestamp: time.Unix(0, int64(span.StartTimeUnixNano)),
	}

	result.Fields, err = OTLPDataToLogRecordFields(span)

	return
}

type ContentType string

const (
	ProtoContentType ContentType = "application/x-protobuf"
	JsonContentType  ContentType = "application/json"
)

func UnmarshalMessage(body []byte, message proto.Message, contentType ContentType) error {
	switch contentType {
	case ProtoContentType:
		err := proto.Unmarshal(body, message)
		if err != nil {
			return err
		}

	case JsonContentType:
		err := json.Unmarshal(body, message)
		if err != nil {
			return err
		}
	}

	return nil
}

type LogRecordsConvertor interface {
	GetMessage() proto.Message
	GetLogRecords() ([]*models.LogRecord, error)
}

type LogsToLogRecordsConvertor struct {
	Message *collectlogs.ExportLogsServiceRequest
}

func (o *LogsToLogRecordsConvertor) GetMessage() proto.Message {
	return o.Message
}

func (o *LogsToLogRecordsConvertor) GetLogRecords() ([]*models.LogRecord, error) {
	logRecords := make([]*models.LogRecord, 0)
	for _, resourceLogs := range o.Message.GetResourceLogs() {
		for _, scopeLogs := range resourceLogs.GetScopeLogs() {
			for _, logRecord := range scopeLogs.GetLogRecords() {
				logRecordModel, err := LogToLogRecord(logRecord)
				if err != nil {
					return nil, err
				}
				logRecords = append(logRecords, logRecordModel)
			}
		}
	}
	return logRecords, nil
}

type TracesToLogRecordsConvertor struct {
	Message *collecttraces.ExportTraceServiceRequest
}

func (o *TracesToLogRecordsConvertor) GetMessage() proto.Message {
	return o.Message
}

func (o *TracesToLogRecordsConvertor) GetLogRecords() ([]*models.LogRecord, error) {
	logRecords := make([]*models.LogRecord, 0)
	for _, resourceLogs := range o.Message.GetResourceSpans() {
		for _, scopeLogs := range resourceLogs.GetScopeSpans() {
			for _, logRecord := range scopeLogs.GetSpans() {
				logRecordModel, err := SpanToLogRecord(logRecord)
				if err != nil {
					return nil, err
				}
				logRecords = append(logRecords, logRecordModel)
			}
		}
	}
	return logRecords, nil
}

type MetricsToLogRecordsConvertor struct {
	Message *collectmetrics.ExportMetricsServiceRequest
}

func (o *MetricsToLogRecordsConvertor) GetMessage() proto.Message {
	return o.Message
}

func (o *MetricsToLogRecordsConvertor) GetLogRecords() ([]*models.LogRecord, error) {
	logRecords := make([]*models.LogRecord, 0)
	for _, resourceLogs := range o.Message.GetResourceMetrics() {
		for _, scopeLogs := range resourceLogs.GetScopeMetrics() {
			for _, logRecord := range scopeLogs.GetMetrics() {
				logRecordModel, err := MetricToLogRecord(logRecord)
				if err != nil {
					return nil, err
				}
				logRecords = append(logRecords, logRecordModel)
			}
		}
	}
	return logRecords, nil
}

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
	Pipeline    string `path:"pipeline" minLength:"1"`
	Body        []byte
	ContentType ContentType
}

func (ior *IngestOTLPRequest) LoadFromHTTPRequest(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	switch contentType := r.Header.Get("Content-Type"); contentType {
	case string(ProtoContentType):
		fallthrough
	case string(JsonContentType):
		ior.ContentType = ContentType(contentType)
		ior.Body = body
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	return nil
}

type IngestOTLPResponse struct {
	Success bool `json:"success"`
}

type OTLPDataType string

const (
	OTLPDataTypeLogs    OTLPDataType = "logs"
	OTLPDataTypeMetrics OTLPDataType = "metrics"
	OTLPDataTypeTraces  OTLPDataType = "traces"
)

func OTLPConvertorFactory(dataType OTLPDataType) (LogRecordsConvertor, error) {
	switch dataType {
	case OTLPDataTypeLogs:
		return &LogsToLogRecordsConvertor{
			Message: &collectlogs.ExportLogsServiceRequest{},
		}, nil
	case OTLPDataTypeMetrics:
		return &MetricsToLogRecordsConvertor{
			Message: &collectmetrics.ExportMetricsServiceRequest{},
		}, nil
	case OTLPDataTypeTraces:
		return &TracesToLogRecordsConvertor{
			Message: &collecttraces.ExportTraceServiceRequest{},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported data type: %s", dataType)
	}
}

func (ctrl *controller) IngestOTLPUsecase(otlDataType OTLPDataType) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req IngestOTLPRequest,
				resp *IngestOTLPResponse,
			) error {
				convertor, err := OTLPConvertorFactory(otlDataType)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to create OTLP convertor",
						slog.String("error", err.Error()),
					)
					return status.Wrap(err, status.Internal)
				}

				message := convertor.GetMessage()
				err = UnmarshalMessage(req.Body, message, req.ContentType)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to parse OTLP message",
						slog.String("error", err.Error()),
					)
					return status.Wrap(err, status.Internal)
				}

				logRecords, err := convertor.GetLogRecords()
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to convert OTLP message to log records",
						slog.String("error", err.Error()),
					)
					return status.Wrap(err, status.Internal)
				}

				err = SendToPipeline(ctx, logRecords, req.Pipeline, ctrl.deps.PipelineRunner, ctrl.deps.ConfigStorage, ctrl.logger)
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

	u.SetName(fmt.Sprintf("ingest_otlp %s", otlDataType))
	u.SetTitle(fmt.Sprintf("Ingest OTLP %s", otlDataType))

	u.SetDescription(fmt.Sprintf("Run otlp %s records through a pipeline", otlDataType))
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
