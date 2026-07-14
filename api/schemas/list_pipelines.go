package schemas

// ListPipelinesRequest is empty: listing pipelines takes no parameters.
type ListPipelinesRequest struct{}

// ListPipelinesResponse carries the names of the available pipelines.
type ListPipelinesResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Pipelines holds the name of every configured pipeline.
	Pipelines []string `json:"pipelines"`
}
