package schemas

// IngestLogsTextRequest carries a plain-text body to push through a pipeline.
type IngestLogsTextRequest struct {
	// Pipeline is the name of the pipeline to run the lines through.
	Pipeline string `path:"pipeline" minLength:"1"`
	// TextBody is the raw text payload; each non-empty line becomes one record.
	TextBody string `contentType:"text/plain"`
}

// IngestLogsTextResponse reports how many lines were processed.
type IngestLogsTextResponse struct {
	// Success reports whether every line was processed.
	Success bool `json:"success"`
	// ProcessedCount is the number of lines that ran through the pipeline.
	ProcessedCount int `json:"processed_count"`
}
