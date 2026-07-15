package forwarders

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"

	"net/http"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	"link-society.com/flowg/internal/models"
)

// datadogRuntime pushes records to the Datadog logs intake API, evaluating the
// source, tags, hostname, message and service dynamic fields per record.
type datadogRuntime struct {
	config *models.ForwarderDatadogV2

	client *http.Client

	ddsourceProg *vm.Program
	ddtagsProg   *vm.Program
	hostnameProg *vm.Program
	messageProg  *vm.Program
	serviceProg  *vm.Program
}

var _ Runtime = (*datadogRuntime)(nil)

func (rt *datadogRuntime) Init(ctx context.Context) error {
	var err error

	rt.client = &http.Client{}

	ddsource := rt.config.DDsource
	if ddsource == "" {
		ddsource = "@expr:log.ddsource"
	}
	rt.ddsourceProg, err = CompileDynamicField(string(ddsource))
	if err != nil {
		return fmt.Errorf("failed to compile ddsource field: %w", err)
	}

	ddtags := rt.config.DDtags
	if ddtags == "" {
		ddtags = "@expr:log.ddtags"
	}
	rt.ddtagsProg, err = CompileDynamicField(string(ddtags))
	if err != nil {
		return fmt.Errorf("failed to compile ddtags field: %w", err)
	}

	hostname := rt.config.Hostname
	if hostname == "" {
		hostname = "@expr:log.hostname"
	}
	rt.hostnameProg, err = CompileDynamicField(string(hostname))
	if err != nil {
		return fmt.Errorf("failed to compile hostname field: %w", err)
	}

	message := rt.config.Message
	if message == "" {
		message = "@expr:log.message"
	}
	rt.messageProg, err = CompileDynamicField(string(message))
	if err != nil {
		return fmt.Errorf("failed to compile message field: %w", err)
	}

	service := rt.config.Service
	if service == "" {
		service = "@expr:log.service"
	}

	rt.serviceProg, err = CompileDynamicField(string(service))
	if err != nil {
		return fmt.Errorf("failed to compile service field: %w", err)
	}

	return nil
}

func (rt *datadogRuntime) Close(context.Context) error {
	return nil
}

func (rt *datadogRuntime) Call(ctx context.Context, record *models.LogRecord) error {
	env := map[string]any{
		"timestamp": record.Timestamp,
		"log":       record.Fields,
	}

	eval := func(prog *vm.Program, field string) (string, error) {
		out, err := expr.Run(prog, env)
		if err != nil {
			return "", fmt.Errorf("failed to evaluate %s expression: %w", field, err)
		}
		str, ok := out.(string)
		if !ok {
			return "", fmt.Errorf("%s expression did not evaluate to string", field)
		}
		return str, nil
	}

	ddsource, err := eval(rt.ddsourceProg, "ddsource")
	if err != nil {
		return fmt.Errorf("failed to evaluate `ddsource` record: %w", err)
	}

	ddtags, err := eval(rt.ddtagsProg, "ddtags")
	if err != nil {
		return fmt.Errorf("failed to evaluate `ddtags` record: %w", err)
	}

	hostname, err := eval(rt.hostnameProg, "hostname")
	if err != nil {
		return fmt.Errorf("failed to evaluate `hostname` record: %w", err)
	}

	message, err := eval(rt.messageProg, "message")
	if err != nil {
		return fmt.Errorf("failed to evaluate `message` record: %w", err)
	}

	service, err := eval(rt.serviceProg, "service")
	if err != nil {
		return fmt.Errorf("failed to evaluate `service` record: %w", err)
	}

	rec := map[string]any{
		"timestamp": record.Timestamp,
		"ddsource":  ddsource,
		"ddtags":    ddtags,
		"hostname":  hostname,
		"message":   message,
		"service":   service,
	}

	payload, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("failed to marshal log record: %w", err)
	}

	buffer := bytes.NewBuffer(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", rt.config.Url, buffer)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("DD-API-KEY", rt.config.ApiKey)

	resp, err := rt.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return fmt.Errorf("can't close body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
