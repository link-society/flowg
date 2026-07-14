package schemas

import "link-society.com/flowg/internal/models"

// GetPipelineRequest identifies the pipeline to retrieve.
type GetPipelineRequest struct {
	// Pipeline is the name of the pipeline to read.
	Pipeline string `path:"pipeline" minLength:"1"`
}

// GetPipelineResponse carries the flow graph of the requested pipeline.
type GetPipelineResponse struct {
	// Success reports whether the pipeline was found and returned.
	Success bool `json:"success"`
	// Flow is the pipeline's flow graph definition.
	Flow *models.FlowGraphV2 `json:"flow"`
}
