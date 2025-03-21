package models

import (
	"fmt"

	"encoding/json"
)

func ConvertFlowGraph(content []byte) (*FlowGraphV2, bool, error) {
	var data map[string]any
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal flow: %w", err)
	}

	version, ok := data["version"].(float64)
	if !ok || version == 0 {
		version = 1
	}

	switch int(version) {
	case 2:
		objV2 := &FlowGraphV2{Version: 2}
		if err := json.Unmarshal(content, objV2); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal flow: %w", err)
		}

		return objV2, false, nil

	case 1:
		objV1 := &FlowGraphV1{Version: 1}
		if err := json.Unmarshal(content, objV1); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal flow: %w", err)
		}

		objV2 := flowGraph_V1_V2(objV1)

		return objV2, true, nil

	default:
		return nil, false, fmt.Errorf("unsupported flow version: %d", int(version))
	}
}

func flowGraph_V1_V2(objV1 *FlowGraphV1) *FlowGraphV2 {
	objV2 := &FlowGraphV2{Version: 2}

	for _, nodeV1 := range objV1.Nodes {
		var nodeV2 *FlowNodeV2

		switch nodeV1.Type {
		case "alert":
			nodeV2 = &FlowNodeV2{
				ID:   nodeV1.ID,
				Type: "forwarder",
				Position: FlowPositionV2{
					X: nodeV1.Position.X,
					Y: nodeV1.Position.Y,
				},
				Data: map[string]string{
					"forwarder": nodeV1.Data["alert"],
				},
			}

		default:
			nodeV2 = &FlowNodeV2{
				ID:   nodeV1.ID,
				Type: nodeV1.Type,
				Position: FlowPositionV2{
					X: nodeV1.Position.X,
					Y: nodeV1.Position.Y,
				},
				Data: nodeV1.Data,
			}
		}

		objV2.Nodes = append(objV2.Nodes, nodeV2)
	}

	for _, edgeV1 := range objV1.Edges {
		edgeV2 := &FlowEdgeV2{
			ID:           edgeV1.ID,
			Source:       edgeV1.Source,
			SourceHandle: edgeV1.SourceHandle,
			Target:       edgeV1.Target,
		}
		objV2.Edges = append(objV2.Edges, edgeV2)
	}

	return objV2
}
