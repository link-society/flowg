package models

import (
	"fmt"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/expr-lang/expr/file"
	"github.com/expr-lang/expr/parser/lexer"
)

// ConvertFlowGraph decodes a stored flow graph of any supported version and
// returns it as the latest FlowGraphV2 (2.1). The boolean reports whether an
// upgrade happened, so callers can persist the migrated form. Supported inputs
// are V1, V2.0 and V2.1.
func ConvertFlowGraph(content []byte) (*FlowGraphV2, bool, error) {
	var data map[string]any
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal flow: %w", err)
	}

	majorVersion, ok := data["version"].(float64)
	if !ok || majorVersion == 0 {
		majorVersion = 1
	}

	minorVersion, ok := data["version.minor"].(float64)
	if !ok {
		minorVersion = 0
	}

	switch int(majorVersion) {
	case 2:
		objV2 := &FlowGraphV2{MajorVersion: 2, MinorVersion: 0}
		if err := json.Unmarshal(content, objV2); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal flow: %w", err)
		}

		switch int(minorVersion) {
		case 1:
			return objV2, false, nil

		case 0:
			objV2m1, err := flowGraphFromV2ToV2m1(objV2)
			if err != nil {
				return nil, false, err
			}

			return objV2m1, true, nil

		default:
			return nil, false, fmt.Errorf("unsupported flow version: %d.%d", int(majorVersion), int(minorVersion))
		}

	case 1:
		objV1 := &FlowGraphV1{Version: 1}
		if err := json.Unmarshal(content, objV1); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal flow: %w", err)
		}

		objV2 := flowGraphFromV1ToV2(objV1)
		objV2m1, err := flowGraphFromV2ToV2m1(objV2)
		if err != nil {
			return nil, false, err
		}

		return objV2m1, true, nil

	default:
		return nil, false, fmt.Errorf("unsupported flow version: %d", int(majorVersion))
	}
}

// flowGraphFromV2ToV2m1 upgrades a 2.0 graph to 2.1 by translating each switch
// node's condition from the legacy filter DSL to expr-lang.
func flowGraphFromV2ToV2m1(objV2 *FlowGraphV2) (*FlowGraphV2, error) {
	objV2.MajorVersion = 2
	objV2.MinorVersion = 1

	for _, nodeV2 := range objV2.Nodes {
		if nodeV2.Type == "switch" {
			if condition, exists := nodeV2.Data["condition"]; exists {
				translated, err := ConvertFilterdslToExprlang(condition)
				if err != nil {
					return nil, err
				}
				nodeV2.Data["condition"] = translated
			}
		}
	}

	return objV2, nil
}

// flowGraphFromV1ToV2 upgrades a V1 graph to V2.0, mapping the old "alert" node
// type onto the "forwarder" node and copying nodes and edges across otherwise
// unchanged.
func flowGraphFromV1ToV2(objV1 *FlowGraphV1) *FlowGraphV2 {
	objV2 := &FlowGraphV2{MajorVersion: 2, MinorVersion: 0}

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

// ConvertFilterdslToExprlang rewrites a legacy filter-DSL expression into an
// equivalent expr-lang one, the only structural change being the single "="
// equality operator becoming "==".
func ConvertFilterdslToExprlang(input string) (string, error) {
	tokens, err := lexer.Lex(file.NewSource(input))
	if err != nil {
		return "", fmt.Errorf("failed to parse expression: %v", err)
	}

	for i, token := range tokens {
		if token.Kind == lexer.Operator && token.Value == "=" {
			tokens[i].Value = "=="
		}
	}

	var values []string
	for _, token := range tokens {
		switch token.Kind {
		case lexer.EOF:
			continue
		case lexer.String:
			values = append(values, strconv.Quote(token.Value))
		default:
			values = append(values, token.Value)
		}
	}
	return strings.Join(values, " "), nil
}
