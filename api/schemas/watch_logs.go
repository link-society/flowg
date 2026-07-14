package schemas

import "io"

// WatchLogsRequest describes a live subscription to a stream's logs.
type WatchLogsRequest struct {
	// Stream is the name of the stream to watch.
	Stream string `path:"stream" minLength:"1"`
	// Filter is an optional filtering expression; only matching records are
	// streamed to the client.
	Filter *string `query:"filter"`
}

// WatchLogsResponse streams matching records to the client as they arrive.
//
// It embeds the writer so the usecase can emit a Server-Sent Events stream
// rather than a single buffered response.
type WatchLogsResponse struct {
	Writer io.Writer
}

func (resp *WatchLogsResponse) SetWriter(w io.Writer) {
	resp.Writer = w
}
