package models

import (
	"fmt"

	"encoding/json"
)

func ConvertWebhook(content []byte) (*WebhookV1, bool, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal webhook: %w", err)
	}

	version, ok := data["version"].(float64)
	if !ok || version == 0 {
		version = 1
	}

	switch int(version) {
	case 1:
		objV1 := &WebhookV1{Version: 1}
		if err := json.Unmarshal(content, objV1); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal webhook: %w", err)
		}

		return objV1, false, nil

	default:
		return nil, false, fmt.Errorf("unsupported webhook version: %d", int(version))
	}
}
