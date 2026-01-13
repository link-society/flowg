package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type forwarderStateDatadogV2 struct {
	client *http.Client

	ddsourceProg *vm.Program
	ddtagsProg   *vm.Program
	hostnameProg *vm.Program
	messageProg  *vm.Program
	serviceProg  *vm.Program
}

type ForwarderDatadogV2 struct {
	Type     string                          `json:"type" enum:"datadog" required:"true"`
	Url      string                          `json:"url" required:"true" format:"uri"`
	ApiKey   string                          `json:"apiKey" required:"true" minLength:"1"`
	DDsource ForwarderDatadogV2DDsourceField `json:"ddsource" required:"true"`
	DDtags   ForwarderDatadogV2DDtagsField   `json:"ddtags" required:"true"`
	Hostname ForwarderDatadogV2HostnameField `json:"hostname" required:"true"`
	Message  ForwarderDatadogV2MessageField  `json:"message" required:"true"`
	Service  ForwarderDatadogV2ServiceField  `json:"service" required:"true"`

	state *forwarderStateDatadogV2
}

func (f *ForwarderDatadogV2) init(ctx context.Context) error {
	var err error
	f.state = &forwarderStateDatadogV2{
		client: &http.Client{},
	}

	ddsource := f.DDsource
	if ddsource == "" {
		ddsource = "@expr:log.ddsource"
	}
	f.state.ddsourceProg, err = CompileDynamicField(string(ddsource))
	if err != nil {
		return fmt.Errorf("failed to compile ddsource field: %w", err)
	}

	ddtags := f.DDtags
	if ddtags == "" {
		ddtags = "@expr:log.ddtags"
	}
	f.state.ddtagsProg, err = CompileDynamicField(string(ddtags))
	if err != nil {
		return fmt.Errorf("failed to compile ddtags field: %w", err)
	}

	hostname := f.Hostname
	if hostname == "" {
		hostname = "@expr:log.hostname"
	}
	f.state.hostnameProg, err = CompileDynamicField(string(hostname))
	if err != nil {
		return fmt.Errorf("failed to compile hostname field: %w", err)
	}

	message := f.Message
	if message == "" {
		message = "@expr:log.message"
	}
	f.state.messageProg, err = CompileDynamicField(string(message))
	if err != nil {
		return fmt.Errorf("failed to compile message field: %w", err)
	}

	service := f.Service
	if service == "" {
		service = "@expr:log.service"
	}

	f.state.serviceProg, err = CompileDynamicField(string(service))
	if err != nil {
		return fmt.Errorf("failed to compile service field: %w", err)
	}

	return nil
}

func (f *ForwarderDatadogV2) close(context.Context) error {
	return nil
}

func (f *ForwarderDatadogV2) call(ctx context.Context, record *LogRecord) error {
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

	ddsource, err := eval(f.state.ddsourceProg, "ddsource")
	if err != nil {
		return fmt.Errorf("failed to evaluate `ddsource` record: %w", err)
	}

	ddtags, err := eval(f.state.ddtagsProg, "ddtags")
	if err != nil {
		return fmt.Errorf("failed to evaluate `ddtags` record: %w", err)
	}

	hostname, err := eval(f.state.hostnameProg, "hostname")
	if err != nil {
		return fmt.Errorf("failed to evaluate `hostname` record: %w", err)
	}

	message, err := eval(f.state.messageProg, "message")
	if err != nil {
		return fmt.Errorf("failed to evaluate `message` record: %w", err)
	}

	service, err := eval(f.state.serviceProg, "service")
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
	req, err := http.NewRequestWithContext(ctx, "POST", f.Url, buffer)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("DD-API-KEY", f.ApiKey)

	resp, err := f.state.client.Do(req)
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
