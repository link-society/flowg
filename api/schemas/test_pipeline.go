package schemas

import "link-society.com/flowg/internal/models"

// TestPipelineRequest carries a pipeline definition and sample records to run
// through it.
type TestPipelineRequest struct {
	// Pipeline is the name used to resolve referenced configuration.
	Pipeline string `path:"pipeline" minLength:"1"`
	// Flow is the flow graph to execute, without persisting it.
	Flow models.FlowGraphV2 `json:"flow" required:"true"`
	// Records are the input log records fed to the pipeline.
	Records []map[string]string `json:"records" required:"true"`
}

// TestPipelineResponse carries the execution trace of the trial run.
type TestPipelineResponse struct {
	// Success reports whether the trial run completed.
	Success bool `json:"success"`
	// Trace records the path each record took through the pipeline nodes.
	Trace []models.PipelineNodeTrace `json:"trace"`
	// Error holds the message of the last record that failed, if any.
	Error *string `json:"error,omitempty"`
}
