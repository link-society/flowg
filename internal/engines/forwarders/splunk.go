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

// splunkRuntime pushes records to a Splunk HTTP Event Collector, evaluating
// the source and host dynamic fields per record.
type splunkRuntime struct {
	config *models.ForwarderSplunkV2

	client *http.Client

	source *vm.Program
	host   *vm.Program
}

var _ Runtime = (*splunkRuntime)(nil)

func (rt *splunkRuntime) Init(context.Context) error {
	var err error

	rt.client = &http.Client{}

	source := rt.config.Source
	if source == "" {
		source = "flowg"
	}
	rt.source, err = CompileDynamicField(string(source))
	if err != nil {
		return fmt.Errorf("failed to compile source field: %w", err)
	}

	host := rt.config.Host
	if host == "" {
		host = "@expr:log.host"
	}
	rt.host, err = CompileDynamicField(string(host))
	if err != nil {
		return fmt.Errorf("failed to compile host field: %w", err)
	}

	return nil
}

func (rt *splunkRuntime) Close(context.Context) error {
	return nil
}

func (rt *splunkRuntime) Call(ctx context.Context, record *models.LogRecord) error {
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

	source, err := eval(rt.source, "source")
	if err != nil {
		return fmt.Errorf("failed to evaluate `source` record: %w", err)
	}

	host, err := eval(rt.host, "host")
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
	req, err := http.NewRequestWithContext(ctx, "POST", rt.config.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add required headers
	req.Header.Add("Authorization", "Splunk "+rt.config.Token)
	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := rt.client.Do(req)
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
