package schemas

import (
	"link-society.com/flowg/internal/models"
)

// UpdateSystemConfigurationRequest is the new system configuration to store. It
// aliases [models.SystemConfiguration] so the whole configuration is replaced
// in one request.
type UpdateSystemConfigurationRequest = models.SystemConfiguration

// UpdateSystemConfigurationResponse reports the outcome of the update.
type UpdateSystemConfigurationResponse = struct {
	// Success reports whether the configuration was persisted.
	Success bool `json:"success"`
}
