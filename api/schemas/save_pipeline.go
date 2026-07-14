package schemas

import "link-society.com/flowg/internal/models"

// SavePipelineRequest carries the pipeline name and its new flow graph.
type SavePipelineRequest struct {
	// Pipeline is the name of the pipeline to create or overwrite.
	Pipeline string `path:"pipeline" minLength:"1"`
	// Flow is the flow graph to store under that name.
	Flow models.FlowGraphV2 `json:"flow" required:"true"`
}

// SavePipelineResponse reports the outcome of the save.
type SavePipelineResponse struct {
	// Success reports whether the pipeline was persisted.
	Success bool `json:"success"`
}
