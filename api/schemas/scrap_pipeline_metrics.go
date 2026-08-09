package schemas

import "io"

// ScrapPipelineMetricsRequest identifies the pipeline to retrieve.
type ScrapPipelineMetricsRequest struct {
	// Pipeline is the name of the pipeline to read.
	Pipeline string `path:"pipeline" minLength:"1"`
}

// ScrapPipelineMetricsResponse embeds the response writer to forward it
// to the pipeline engine
type ScrapPipelineMetricsResponse struct {
	Writer io.Writer
}

func (resp *ScrapPipelineMetricsResponse) SetWriter(w io.Writer) {
	resp.Writer = w
}
