package schemas

import (
	"mime/multipart"
)

// RestoreAuthRequest carries the authentication database snapshot to load.
type RestoreAuthRequest struct {
	// Backup is the uploaded snapshot, as produced by
	// [NewBackupAuthUsecase].
	Backup multipart.File `formData:"backup"`
}

// RestoreAuthResponse reports the outcome of the restore.
type RestoreAuthResponse struct {
	// Success reports whether the snapshot was loaded.
	Success bool `json:"success"`
}
