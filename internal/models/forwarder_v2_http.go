package models

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"

	"net/http"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type forwarderStateHttpV2 struct {
	client *http.Client

	body *vm.Program
}

type ForwarderHttpV2 struct {
	Type    string                   `json:"type" enum:"http" required:"true"`
	Url     string                   `json:"url" required:"true" format:"uri"`
	Headers map[string]string        `json:"headers,omitempty"`
	Body    ForwarderHttpV2BodyField `json:"body,omitempty"`

	state *forwarderStateHttpV2
}

func (f *ForwarderHttpV2) init(context.Context) error {
	var err error
	f.state = &forwarderStateHttpV2{
		client: &http.Client{},
	}

	body := f.Body
	if body == "" {
		body = "@expr:toJSON(log)"
	}
	f.state.body, err = CompileDynamicField(string(body))
	if err != nil {
		return fmt.Errorf("failed to compile body field: %w", err)
	}

	return nil
}

func (f *ForwarderHttpV2) close(context.Context) error {
	return nil
}

func (f *ForwarderHttpV2) call(ctx context.Context, record *LogRecord) error {
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

	body, err := eval(f.state.body, "body")
	if err != nil {
		return fmt.Errorf("failed to evaluate `body` record: %w", err)
	}

	rec := map[string]any{
		"timestamp": record.Timestamp,
		"body":      body,
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

	for key, value := range f.Headers {
		req.Header.Add(key, value)
	}

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
