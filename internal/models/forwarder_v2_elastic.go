package models

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v9"
)

type forwarderStateElasticV2 struct {
	client *elasticsearch.Client
}

type ForwarderElasticV2 struct {
	Type      string   `json:"type" enum:"elastic" required:"true"`
	Index     string   `json:"index" required:"true"`
	Username  string   `json:"username" required:"true"`
	Password  string   `json:"password" required:"true"`
	Addresses []string `json:"addresses" required:"true"`
	CACert    string   `json:"ca,omitempty"`

	state *forwarderStateElasticV2
}

func (f *ForwarderElasticV2) init(context.Context) error {
	var caBytes []byte
	if f.CACert != "" {
		caBytes = []byte(f.CACert)
	}

	cfg := elasticsearch.Config{
		Username:  f.Username,
		Password:  f.Password,
		Addresses: f.Addresses,
		CACert:    caBytes,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create ElasticSearch client: %w", err)
	}

	f.state = &forwarderStateElasticV2{
		client: client,
	}

	return nil
}

func (f *ForwarderElasticV2) close(context.Context) error {
	return nil
}

func (f *ForwarderElasticV2) call(ctx context.Context, record *LogRecord) error {
	resp, err := f.state.client.Indices.Exists([]string{f.Index}, f.state.client.Indices.Exists.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to send ElasticSearch request: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		if resp.StatusCode == 404 {
			resp, err := f.state.client.Indices.Create(f.Index, f.state.client.Indices.Create.WithContext(ctx))
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

	resp, err = f.state.client.Index(f.Index, data, f.state.client.Index.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to send ElasticSearch request: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("failed to index log record: %s", resp.String())
	}

	return nil
}
