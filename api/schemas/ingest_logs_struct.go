package schemas

// IngestLogsStructRequest carries structured records to push through a pipeline.
type IngestLogsStructRequest struct {
	// Pipeline is the name of the pipeline to run the records through.
	Pipeline string `path:"pipeline" minLength:"1"`
	// Records are the structured log records to ingest.
	Records []map[string]string `json:"records" required:"true"`
}

// IngestLogsStructResponse reports how many records were processed.
type IngestLogsStructResponse struct {
	// Success reports whether every record was processed.
	Success bool `json:"success"`
	// ProcessedCount is the number of records that ran through the pipeline.
	ProcessedCount int `json:"processed_count"`
}
