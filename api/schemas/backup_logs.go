package schemas

import "io"

// BackupLogsRequest is empty: the whole log database is exported.
type BackupLogsRequest struct{}

// BackupLogsResponse streams the log database snapshot to the client.
//
// It embeds the writer so the snapshot is streamed as a file download rather
// than buffered in memory.
type BackupLogsResponse struct {
	Writer io.Writer
}

func (resp *BackupLogsResponse) SetWriter(w io.Writer) {
	resp.Writer = w
}
