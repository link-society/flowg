package schemas

// TestTransformerRequest carries the transformer source and a sample record to
// run it against.
type TestTransformerRequest struct {
	// Code is the VRL source to evaluate, without persisting it.
	Code string `json:"code" required:"true"`
	// Record is the input log record fed to the transformer.
	Record map[string]string `json:"record" required:"true"`
}

// TestTransformerResponse carries the records produced by the trial run.
type TestTransformerResponse struct {
	// Success reports whether the script compiled and ran.
	Success bool `json:"success"`
	// Records holds the output records emitted by the transformer.
	Records []map[string]string `json:"records"`
}
