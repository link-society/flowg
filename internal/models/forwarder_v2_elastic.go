package models

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v9"
)

type ForwarderElasticV2 struct {
	Type      string   `json:"type" enum:"elastic"`
	Index     string   `json:"index"`
	Addresses []string `json:"addresses"`
	CACert    string   `json:"ca,omitempty"`
	Token     string   `json:"token,omitempty"`
}

func (f *ForwarderElasticV2) call(ctx context.Context, record *LogRecord) error {
	cfg := elasticsearch.Config{
		CACert:    []byte(f.CACert),
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create ElasticSearch client: %w", err)
	}

	resp, err := client.Indices.Exists([]string{f.Index}, client.Indices.Exists.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to send ElasticSearch request: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		if resp.StatusCode == 404 {
			resp, err := client.Indices.Create(f.Index, client.Indices.Create.WithContext(ctx))
			if err != nil {
				return fmt.Errorf("failed to send ElasticSearch request: %w", err)
			}
			defer resp.Body.Close()

			if resp.IsError() {
				return fmt.Errorf("failed to check index '%s': %s", f.Index, resp.String())
			}
		} else {
			return fmt.Errorf("failed to check index '%s': %s", f.Index, resp.String())
		}
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal log record: %w", err)
	}
	data := bytes.NewReader(payload)

	resp, err = client.Index(f.Index, data, client.Index.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to send ElasticSearch request: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("failed to index log record: %s", resp.String())
	}

	return nil
}
