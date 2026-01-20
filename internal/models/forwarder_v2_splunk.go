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

type forwarderStateSplunkV2 struct {
	client *http.Client

	source *vm.Program
	host   *vm.Program
}

type ForwarderSplunkV2 struct {
	Type     string                       `json:"type" enum:"splunk" required:"true"`
	Endpoint string                       `json:"endpoint" required:"true" format:"uri"`
	Token    string                       `json:"token" required:"true" minLength:"1"`
	Source   ForwarderSplunkV2SourceField `json:"source"`
	Host     ForwarderSplunkV2HostField   `json:"host"`

	state *forwarderStateSplunkV2
}

func (f *ForwarderSplunkV2) init(context.Context) error {
	var err error
	f.state = &forwarderStateSplunkV2{
		client: &http.Client{},
	}

	source := f.Source
	if source == "" {
		source = "flowg"
	}
	f.state.source, err = CompileDynamicField(string(source))
	if err != nil {
		return fmt.Errorf("failed to compile source field: %w", err)
	}

	host := f.Host
	if host == "" {
		host = "@expr:log.host"
	}
	f.state.host, err = CompileDynamicField(string(host))
	if err != nil {
		return fmt.Errorf("failed to compile host field: %w", err)
	}

	return nil
}

func (f *ForwarderSplunkV2) close(context.Context) error {
	return nil
}

func (f *ForwarderSplunkV2) call(ctx context.Context, record *LogRecord) error {
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

	source, err := eval(f.state.source, "source")
	if err != nil {
		return fmt.Errorf("failed to evaluate `source` record: %w", err)
	}

	host, err := eval(f.state.host, "host")
	if err != nil {
		return fmt.Errorf("failed to evaluate `host` record: %w", err)
	}

	// Convert map[string]string to map[string]interface{}
	eventFields := make(map[string]interface{})
	for k, v := range record.Fields {
		eventFields[k] = v
	}

	// Create Splunk HEC payload
	payload := struct {
		Event      map[string]interface{} `json:"event"`
		Sourcetype string                 `json:"sourcetype"`
		Source     string                 `json:"source"`
		Host       string                 `json:"host"`
		Time       int64                  `json:"time"`
	}{
		Event:      eventFields,
		Sourcetype: "json",
		Source:     source,
		Host:       host,
		Time:       record.Timestamp.Unix(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal log record: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", f.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add required headers
	req.Header.Add("Authorization", "Splunk "+f.Token)
	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := f.state.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from Splunk: %d", resp.StatusCode)
	}

	var result struct {
		Text string `json:"text"`
		Code int    `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	_ = resp.Body.Close()
	return nil
}
