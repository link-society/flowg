package schemas

// DeletePipelineRequest identifies the pipeline to remove.
type DeletePipelineRequest struct {
	// Pipeline is the name of the pipeline to delete.
	Pipeline string `path:"pipeline" minLength:"1"`
}

// DeletePipelineResponse reports the outcome of the deletion.
type DeletePipelineResponse struct {
	// Success reports whether the pipeline was removed.
	Success bool `json:"success"`
}
