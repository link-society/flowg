package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ForwarderDatadogV2 struct {
	Type    string `json:"type" enum:"datadog"`
	Url     string `json:"url"`
	ApiKey  string `json:"apiKey"`
	Source  string `json:"source"`
	Service string `json:"service"`
}

type DatadogLogItem struct {
	DdSource string `json:"ddsource"`
	DdTags   string `json:"ddtags"`
	Hostname string `json:"hostname"`
	Message  string `json:"message"`
	Service  string `json:"service"`
}

func (f *ForwarderDatadogV2) call(ctx context.Context, record *LogRecord) error {
	logItems := []*DatadogLogItem{CreateDatadogHttpLogItem(f, record)}

	payload, err := json.Marshal(logItems)
	if err != nil {
		return fmt.Errorf("failed to marshal log record: %w", err)
	}

	buffer := bytes.NewBuffer(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", f.Url, buffer)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	//TODO is this hardcoded the way to go?
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("DD-API-KEY", f.ApiKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func CreateDatadogHttpLogItem(f *ForwarderDatadogV2, record *LogRecord) *DatadogLogItem {
	logItem := &DatadogLogItem{}
	var tags []string

	for key, value := range record.Fields {
		switch key {
		case f.Service:
			logItem.Service = value
		case f.Source:
			logItem.DdSource = value
		case "message": /* TODO we should not hard code this */
			logItem.Message = value
		case "hostname": /* TODO we should not hard code this */
			logItem.Hostname = value
		default:
			tags = append(tags, fmt.Sprintf("%s:%s", key, value))
		}
	}
	logItem.DdTags = strings.Join(tags, ",")
	return logItem
}
