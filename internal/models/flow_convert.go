package models

import (
	"fmt"

	"encoding/json"
)

func ConvertFlowGraph(content []byte) (*FlowGraphV1, bool, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal flow: %w", err)
	}

	version, ok := data["version"].(float64)
	if !ok || version == 0 {
		version = 1
	}

	switch int(version) {
	case 1:
		objV1 := &FlowGraphV1{Version: 1}
		if err := json.Unmarshal(content, objV1); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal flow: %w", err)
		}

		return objV1, false, nil

	default:
		return nil, false, fmt.Errorf("unsupported flow version: %d", int(version))
	}
}
