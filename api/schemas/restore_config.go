package schemas

import (
	"mime/multipart"
)

// RestoreConfigRequest carries the configuration database snapshot to load.
type RestoreConfigRequest struct {
	// Backup is the uploaded snapshot, as produced by the backup-config
	// endpoint (see [BackupConfigResponse]).
	Backup multipart.File `formData:"backup"`
}

// RestoreConfigResponse reports the outcome of the restore.
type RestoreConfigResponse struct {
	// Success reports whether the snapshot was loaded.
	Success bool `json:"success"`
}
