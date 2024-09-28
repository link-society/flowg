package pipelines

import "fmt"

type InvalidFlowNodeTypeError struct {
	Type string
}

func (e *InvalidFlowNodeTypeError) Error() string {
	return fmt.Sprintf("unknown node type: %s", e.Type)
}

type MissingFlowNodeDataError struct {
	NodeID string
	Key    string
}

func (e *MissingFlowNodeDataError) Error() string {
	return fmt.Sprintf("missing data key %s for node %s", e.Key, e.NodeID)
}

type InvalidFlowEdgeError struct {
	Source string
	Target string
}

func (e *InvalidFlowEdgeError) Error() string {
	return fmt.Sprintf("invalid edge: %s -> %s", e.Source, e.Target)
}

type InvalidEntrypointError struct {
	Entrypoint string
}

func (e *InvalidEntrypointError) Error() string {
	return fmt.Sprintf("invalid entrypoint: %s", e.Entrypoint)
}
