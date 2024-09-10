package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"link-society.com/flowg/internal/data/alerting"
)

const ALERTS_STORAGE_TYPE = "alerts"

type AlertSystem struct {
	backend *Storage
}

func NewAlertSystem(backend *Storage) *AlertSystem {
	backend.resolveStorageTypeDir(ALERTS_STORAGE_TYPE)
	return &AlertSystem{backend: backend}
}

func (sys *AlertSystem) List() ([]string, error) {
	items, err := sys.backend.listStorageTypeItems(ALERTS_STORAGE_TYPE)
	if err != nil {
		return nil, err
	}

	results := []string{}

	for _, item := range items {
		if strings.HasSuffix(item, ".json.b64") {
			results = append(results, item[:len(item)-9])
		}
	}

	return results, nil
}

func (sys *AlertSystem) Read(name string) (*alerting.Webhook, error) {
	b64content, err := sys.backend.readStorageTypeItem(ALERTS_STORAGE_TYPE, name+".json.b64")
	if err != nil {
		return nil, err
	}

	content, err := base64.StdEncoding.DecodeString(b64content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode webhook %s: %w", name, err)
	}

	webhook := &alerting.Webhook{}
	err = json.Unmarshal(content, webhook)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook %s: %w", name, err)
	}

	return webhook, nil
}

func (sys *AlertSystem) Write(name string, webhook *alerting.Webhook) error {
	content, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook %s: %w", name, err)
	}

	b64content := base64.StdEncoding.EncodeToString(content)

	return sys.backend.writeStorageTypeItem(ALERTS_STORAGE_TYPE, name+".json.b64", b64content)
}

func (sys *AlertSystem) Delete(name string) error {
	return sys.backend.deleteStorageTypeItem(ALERTS_STORAGE_TYPE, name+".json.b64")
}
