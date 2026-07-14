package schemas

import "io"

// BackupAuthRequest is empty: the whole authentication database is exported.
type BackupAuthRequest struct{}

// BackupAuthResponse streams the authentication database snapshot to the client.
//
// It embeds the writer so the snapshot is streamed as a file download rather
// than buffered in memory.
type BackupAuthResponse struct {
	Writer io.Writer
}

func (resp *BackupAuthResponse) SetWriter(w io.Writer) {
	resp.Writer = w
}
