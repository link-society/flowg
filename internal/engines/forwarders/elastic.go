package forwarders

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v9"

	"link-society.com/flowg/internal/models"
)

// elasticRuntime indexes records into an Elasticsearch index, creating the
// index on first delivery if it does not exist.
type elasticRuntime struct {
	config *models.ForwarderElasticV2

	client *elasticsearch.Client
}

var _ Runtime = (*elasticRuntime)(nil)

func (rt *elasticRuntime) Init(context.Context) error {
	opts := []elasticsearch.Option{
		elasticsearch.WithAddresses(rt.config.Addresses...),
		elasticsearch.WithBasicAuth(rt.config.Username, rt.config.Password),
	}

	if rt.config.CACert != "" {
		caBytes := []byte(rt.config.CACert)
		opts = append(opts, elasticsearch.WithCACert(caBytes))
	}

	client, err := elasticsearch.New(opts...)
	if err != nil {
		return fmt.Errorf("failed to create ElasticSearch client: %w", err)
	}

	rt.client = client

	return nil
}

func (rt *elasticRuntime) Close(context.Context) error {
	return nil
}

func (rt *elasticRuntime) Call(ctx context.Context, record *models.LogRecord) error {
	resp, err := rt.client.Indices.Exists([]string{rt.config.Index}, rt.client.Indices.Exists.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to send ElasticSearch request: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		if resp.StatusCode == 404 {
			resp, err := rt.client.Indices.Create(rt.config.Index, rt.client.Indices.Create.WithContext(ctx))
			if err != nil {
				return fmt.Errorf("failed to send ElasticSearch request: %w", err)
			}
			defer resp.Body.Close()

			if resp.IsError() {
				return fmt.Errorf("failed to check index '%s': %s", rt.config.Index, resp.String())
			}
		} else {
			return fmt.Errorf("failed to check index '%s': %s", rt.config.Index, resp.String())
		}
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal log record: %w", err)
	}
	data := bytes.NewReader(payload)

	resp, err = rt.client.Index(rt.config.Index, data, rt.client.Index.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to send ElasticSearch request: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("failed to index log record: %s", resp.String())
	}

	return nil
}
