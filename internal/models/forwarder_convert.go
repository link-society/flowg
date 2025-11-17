package models

import (
	"fmt"

	"encoding/json"
)

func ConvertForwarder(content []byte) (*ForwarderV2, bool, error) {
	var data map[string]any
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal webhook: %w", err)
	}

	version, ok := data["version"].(float64)
	if !ok || version == 0 {
		version = 1
	}

	switch int(version) {
	case 2:
		objV2 := &ForwarderV2{Version: 2}
		if err := json.Unmarshal(content, objV2); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal webhook: %w", err)
		}

		return objV2, false, nil

	case 1:
		objV1 := &ForwarderV1{Version: 1}
		if err := json.Unmarshal(content, objV1); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal forwarder: %w", err)
		}

		objV2 := forwarderFromV1ToV2(objV1)

		return objV2, true, nil

	default:
		return nil, false, fmt.Errorf("unsupported forwarder version: %d", int(version))
	}
}

func forwarderFromV1ToV2(objV1 *ForwarderV1) *ForwarderV2 {
	return &ForwarderV2{
		Version: 2,
		Config: ForwarderConfigV2{
			Http: &ForwarderHttpV2{
				Type:    "http",
				Url:     objV1.Url,
				Headers: objV1.Headers,
			},
		},
	}
}
