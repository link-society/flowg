package pipelines

import (
	"context"
	"errors"

	"link-society.com/flowg/internal/app/metrics"

	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// Pipeline is a compiled, ready-to-run flow: a graph of Nodes indexed by ID,
// with the subset of source nodes that act as named entrypoints.
type Pipeline struct {
	Name        string
	Entrypoints map[string]Node
	nodes       map[string]Node
}

// BuildFromStorage loads the persisted flow graph for name and compiles it into
// a runnable Pipeline.
func BuildFromStorage(ctx context.Context, configStorage storage.ConfigStorage, name string) (*Pipeline, error) {
	flowGraph, err := configStorage.ReadPipeline(ctx, name)
	if err != nil {
		return nil, err
	}

	return BuildFlow(ctx, configStorage, name, flowGraph)
}

// BuildFlow compiles a flow graph into a Pipeline: it turns each flow node into
// the matching Node implementation (resolving referenced transformers and
// forwarders from storage), then wires the edges as each node's successors.
// Source nodes become entrypoints keyed by their declared type.
func BuildFlow(ctx context.Context, configStorage storage.ConfigStorage, name string, flowGraph *models.FlowGraphV2) (*Pipeline, error) {
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
				ID:          flowNode.ID,
				Transformer: transformer,
			}
			pipelineNodes[flowNode.ID] = pipelineNode

		case "switch":
			condition, exists := flowNode.Data["condition"]
			if !exists {
				return nil, &MissingFlowNodeDataError{
					NodeID: flowNode.ID,
					Key:    "condition",
				}
			}

			pipelineNode := &SwitchNode{
				ID:        flowNode.ID,
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
				ID:       flowNode.ID,
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
				ID:        flowNode.ID,
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
				ID:     flowNode.ID,
				Stream: stream,
			}
			pipelineNodes[flowNode.ID] = pipelineNode

		case "source":
			sourceType, exists := flowNode.Data["type"]
			if !exists {
				sourceType = DIRECT_ENTRYPOINT
			}

			pipelineNode := &SourceNode{
				ID: flowNode.ID,
			}
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

// Init initialises every node (compiling VRL transformers, opening forwarder
// connections, ...), joining all errors so callers see every failure.
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

// Close releases every node's resources, joining all errors.
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

// Process feeds a record into the node selected by entrypoint and records the
// pipeline-level success metric.
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
