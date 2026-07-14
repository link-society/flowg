package forwarders

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"

	"net/http"
	"net/url"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	"link-society.com/flowg/internal/models"
)

type httpRuntime struct {
	config *models.ForwarderHttpV2

	client *http.Client

	body *vm.Program
}

var _ Runtime = (*httpRuntime)(nil)

func (rt *httpRuntime) Init(ctx context.Context) error {
	var err error
	rt.client = &http.Client{}

	body := rt.config.Body
	if body == "" {
		body = "@expr:toJSON(log)"
	}
	rt.body, err = CompileDynamicField(string(body))
	if err != nil {
		return fmt.Errorf("failed to compile body field: %w", err)
	}

	return nil
}

func (rt *httpRuntime) Close(ctx context.Context) error {
	return nil
}

func (rt *httpRuntime) Call(ctx context.Context, record *models.LogRecord) error {
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

	body, err := eval(rt.body, "body")
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
	req, err := http.NewRequestWithContext(ctx, "POST", rt.config.Url, buffer)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range rt.config.Headers {
		req.Header.Add(key, value)
	}

	if len(rt.config.Proxy) > 0 {
		proxy, err := url.Parse(rt.config.Proxy)
		if err != nil {
			return fmt.Errorf("failed to parse proxy: %w", err)
		}
		rt.client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	}

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
