package schemas

import "io"

// BackupConfigRequest is empty: the whole configuration database is exported.
type BackupConfigRequest struct{}

// BackupConfigResponse streams the configuration database snapshot to the
// client.
//
// It embeds the writer so the snapshot is streamed as a file download rather
// than buffered in memory.
type BackupConfigResponse struct {
	Writer io.Writer
}

func (resp *BackupConfigResponse) SetWriter(w io.Writer) {
	resp.Writer = w
}
