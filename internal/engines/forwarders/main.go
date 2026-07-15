package forwarders

import (
	"context"
	"errors"

	"strconv"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	"link-society.com/flowg/internal/models"
)

// ErrNotImplemented is returned by NewRuntime when the forwarder configuration
// selects no known backend.
var ErrNotImplemented = errors.New("runtime not implemented")

// Runtime executes a forwarder: it delivers log records to the external
// destination described by a forwarder configuration.
//
// Init compiles the configuration's dynamic fields and builds the backend
// client; Call delivers one record; Close releases the connection for the
// backends that hold one.
type Runtime interface {
	Init(ctx context.Context) error
	Close(ctx context.Context) error
	Call(ctx context.Context, record *models.LogRecord) error
}

// NewRuntime returns the Runtime implementation matching the configuration's
// tagged union, or ErrNotImplemented when no backend is selected.
func NewRuntime(cfg *models.ForwarderV2) (Runtime, error) {
	switch {
	case cfg.Config.Http != nil:
		return &httpRuntime{config: cfg.Config.Http}, nil

	case cfg.Config.Syslog != nil:
		return &syslogRuntime{config: cfg.Config.Syslog}, nil

	case cfg.Config.Datadog != nil:
		return &datadogRuntime{config: cfg.Config.Datadog}, nil

	case cfg.Config.Amqp != nil:
		return &amqpRuntime{config: cfg.Config.Amqp}, nil

	case cfg.Config.Splunk != nil:
		return &splunkRuntime{config: cfg.Config.Splunk}, nil

	case cfg.Config.Otlp != nil:
		return &otlpRuntime{config: cfg.Config.Otlp}, nil

	case cfg.Config.Elastic != nil:
		return &elasticRuntime{config: cfg.Config.Elastic}, nil

	case cfg.Config.Clickhouse != nil:
		return &clickhouseRuntime{config: cfg.Config.Clickhouse}, nil

	case cfg.Config.AwsCloudWatch != nil:
		return &awsCloudWatchRuntime{config: cfg.Config.AwsCloudWatch}, nil

	case cfg.Config.GoogleCloudLogging != nil:
		return &googleCloudLoggingRuntime{config: cfg.Config.GoogleCloudLogging}, nil

	case cfg.Config.AzureMonitor != nil:
		return &azureMonitorRuntime{config: cfg.Config.AzureMonitor}, nil

	default:
		return nil, ErrNotImplemented
	}
}

// CompileDynamicField compiles a dynamic field into an expr program. A value
// prefixed with "@expr:" is compiled as an expression; any other value is
// compiled as a quoted string literal, so plain values evaluate to themselves.
func CompileDynamicField(value string) (*vm.Program, error) {
	if len(value) >= 6 && value[:6] == "@expr:" {
		return expr.Compile(
			value[6:],
			expr.Env(map[string]any{}),
			expr.AllowUndefinedVariables(),
		)
	} else {
		return expr.Compile(
			strconv.Quote(value),
			expr.Env(map[string]any{}),
			expr.AllowUndefinedVariables(),
		)
	}
}
