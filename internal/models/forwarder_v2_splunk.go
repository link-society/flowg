package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ForwarderSplunkV2 struct {
	Type     string `json:"type" enum:"splunk"`
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
}

func (f *ForwarderSplunkV2) call(ctx context.Context, record *LogRecord) error {
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
		Source:     "flowg",
		Host:       getHost(record.Fields),
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
	client := &http.Client{
		Timeout: 500 * time.Millisecond, // Increased timeout for better reliability
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Splunk: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from Splunk: %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		Text string `json:"text"`
		Code int    `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode Splunk response: %w", err)
	}

	// Check response code
	if result.Code != 0 {
		return fmt.Errorf("Splunk returned error: %s", result.Text)
	}

	return nil
}

// getHost returns the host from fields or a default value
func getHost(fields map[string]string) string {
	host, ok := fields["host"]
	if !ok || host == "" {
		return "flowg"
	}
	return host
}
