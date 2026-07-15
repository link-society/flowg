package models

// PipelineNodeTrace is the record of a single node's execution during a
// pipeline dry run: its input fields, the records it emitted, and any error.
type PipelineNodeTrace struct {
	NodeID string              `json:"nodeID"`
	Input  map[string]string   `json:"input"`
	Output []map[string]string `json:"output"`
	Error  *string             `json:"error"`
}
