package schemas

import (
	"mime/multipart"
)

// RestoreLogsRequest carries the log database snapshot to load.
type RestoreLogsRequest struct {
	// Backup is the uploaded snapshot, as produced by
	// [NewBackupLogsUsecase].
	Backup multipart.File `formData:"backup"`
}

// RestoreLogsResponse reports the outcome of the restore.
type RestoreLogsResponse struct {
	// Success reports whether the snapshot was loaded.
	Success bool `json:"success"`
}
