package pipelines

import (
	"context"
	"errors"

	"link-society.com/flowg/internal/app/metrics"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/config"

	"link-society.com/flowg/internal/utils/ffi/filterdsl"
)

type Pipeline struct {
	Name        string
	Entrypoints map[string]Node
	nodes       map[string]Node
}

func Build(ctx context.Context, configStorage config.Storage, name string) (*Pipeline, error) {
	flowGraph, err := configStorage.ReadPipeline(ctx, name)
	if err != nil {
		return nil, err
	}

	var (
		pipelineNodes   = make(map[string]Node)
		flowNodesByID   = make(map[string]*models.FlowNodeV2)
		sourceNodeTypes = make(map[string]string)
		entrypointNodes = make(map[string]Node)
	)

	for _, flowNode := range flowGraph.Nodes {
		flowNodesByID[flowNode.ID] = flowNode

		switch flowNode.Type {
		case "transform":
			transformerName, exists := flowNode.Data["transformer"]
			if !exists {
				return nil, &MissingFlowNodeDataError{
					NodeID: flowNode.ID,
					Key:    "transformer",
				}
			}

			transformer, err := configStorage.ReadTransformer(ctx, transformerName)
			if err != nil {
				return nil, err
			}

			pipelineNode := &TransformNode{
				Transformer: transformer,
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

		case "pipeline":
			pipeline, exists := flowNode.Data["pipeline"]
			if !exists {
				return nil, &MissingFlowNodeDataError{
					NodeID: flowNode.ID,
					Key:    "pipeline",
				}
			}

			pipelineNode := &PipelineNode{
				Pipeline: pipeline,
			}
			pipelineNodes[flowNode.ID] = pipelineNode

		case "forwarder":
			forwarderName, exists := flowNode.Data["forwarder"]
			if !exists {
				return nil, &MissingFlowNodeDataError{
					NodeID: flowNode.ID,
					Key:    "forwarder",
				}
			}

			forwarder, err := configStorage.ReadForwarder(ctx, forwarderName)
			if err != nil {
				return nil, err
			}

			pipelineNode := &ForwardNode{
				Forwarder: forwarder,
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

		case "source":
			sourceType, exists := flowNode.Data["type"]
			if !exists {
				sourceType = DIRECT_ENTRYPOINT
			}

			pipelineNode := &SourceNode{}
			pipelineNodes[flowNode.ID] = pipelineNode
			sourceNodeTypes[flowNode.ID] = sourceType
			entrypointNodes[sourceType] = pipelineNode

		default:
			return nil, &InvalidFlowNodeTypeError{Type: flowNode.Type}
		}
	}

	for _, flowEdge := range flowGraph.Edges {
		sourceNode, sourceExists := pipelineNodes[flowEdge.Source]
		targetNode, targetExists := pipelineNodes[flowEdge.Target]

		if !sourceExists || !targetExists {
			return nil, &InvalidFlowEdgeError{
				Source: flowEdge.Source,
				Target: flowEdge.Target,
			}
		}

		switch source := sourceNode.(type) {
		case *SourceNode:
			source.Next = append(source.Next, targetNode)

		case *TransformNode:
			source.Next = append(source.Next, targetNode)

		case *SwitchNode:
			source.Next = append(source.Next, targetNode)

		default:
			panic("unreachable")
		}
	}

	return &Pipeline{
		Name:        name,
		Entrypoints: entrypointNodes,
		nodes:       pipelineNodes,
	}, nil
}

func (p *Pipeline) Init(ctx context.Context) error {
	var errs []error

	for _, node := range p.nodes {
		if err := node.Init(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (p *Pipeline) Close(ctx context.Context) error {
	var errs []error

	for _, node := range p.nodes {
		if err := node.Close(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (p *Pipeline) Process(
	ctx context.Context,
	entrypoint string,
	record *models.LogRecord,
) error {
	rootNode, exists := p.Entrypoints[entrypoint]
	if !exists {
		return &InvalidEntrypointError{Entrypoint: entrypoint}
	}

	err := rootNode.Process(ctx, record)
	metrics.IncPipelineLogCounter(p.Name, err == nil)
	return err
}
