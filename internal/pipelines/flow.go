package pipelines

import "link-society.com/flowg/internal/filterdsl"

type FlowGraph struct {
	Nodes []*FlowNode `json:"nodes"`
	Edges []*FlowEdge `json:"edges"`
}

type FlowNode struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Position FlowPosition      `json:"position"`
	Data     map[string]string `json:"data"`
}

type FlowEdge struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle"`
	Target       string `json:"target"`
}

type FlowPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (flowGraph FlowGraph) BuildPipeline() (*Pipeline, error) {
	pipelineNodes := make(map[string]Node)
	flowNodesByID := make(map[string]*FlowNode)

	rootFlowNodeId := ""
	var rootPipelineNode Node

	for _, flowNode := range flowGraph.Nodes {
		flowNodesByID[flowNode.ID] = flowNode

		switch flowNode.Type {
		case "transform":
			transformer, exists := flowNode.Data["transformer"]
			if !exists {
				return nil, &MissingFlowNodeDataError{
					NodeID: flowNode.ID,
					Key:    "transformer",
				}
			}

			pipelineNode := &TransformNode{
				TransformerName: transformer,
			}
			pipelineNodes[flowNode.ID] = pipelineNode

		case "switch":
			conditionSource, exists := flowNode.Data["condition"]
			if !exists {
				return nil, &MissingFlowNodeDataError{
					NodeID: flowNode.ID,
					Key:    "condition",
				}
			}

			condition, err := filterdsl.Compile(conditionSource)
			if err != nil {
				return nil, err
			}

			pipelineNode := &SwitchNode{
				Condition: condition,
			}
			pipelineNodes[flowNode.ID] = pipelineNode

		case "router":
			stream, exists := flowNode.Data["stream"]
			if !exists {
				return nil, &MissingFlowNodeDataError{
					NodeID: flowNode.ID,
					Key:    "stream",
				}
			}

			pipelineNode := &RouterNode{
				Stream: stream,
			}
			pipelineNodes[flowNode.ID] = pipelineNode

		case "root":
			rootFlowNodeId = flowNode.ID

		default:
			return nil, &InvalidFlowNodeTypeError{Type: flowNode.Type}
		}
	}

	for _, flowEdge := range flowGraph.Edges {
		if flowEdge.Source == rootFlowNodeId {
			targetNode, targetExists := pipelineNodes[flowEdge.Target]
			if !targetExists {
				return nil, &InvalidFlowEdgeError{
					Source: flowEdge.Source,
					Target: flowEdge.Target,
				}
			}

			rootPipelineNode = targetNode
		} else {
			sourceNode, sourceExists := pipelineNodes[flowEdge.Source]
			targetNode, targetExists := pipelineNodes[flowEdge.Target]

			if !sourceExists || !targetExists {
				return nil, &InvalidFlowEdgeError{
					Source: flowEdge.Source,
					Target: flowEdge.Target,
				}
			}

			switch source := sourceNode.(type) {
			case *TransformNode:
				source.Next = append(source.Next, targetNode)

			case *SwitchNode:
				source.Next = append(source.Next, targetNode)

			default:
				panic("unreachable")
			}
		}
	}

	if rootPipelineNode == nil {
		return nil, &MissingFlowRootNodeError{}
	}

	return &Pipeline{
		Root: rootPipelineNode,
	}, nil
}
