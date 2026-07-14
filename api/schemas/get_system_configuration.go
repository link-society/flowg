package schemas

import (
	"link-society.com/flowg/internal/models"
)

// GetSystemConfigurationRequest is empty: the system configuration is global.
type GetSystemConfigurationRequest struct{}

// GetSystemConfigurationResponse carries the current system configuration.
type GetSystemConfigurationResponse = struct {
	// Success reports whether the configuration was returned.
	Success bool `json:"success"`
	// Configuration is the current global system configuration.
	Configuration models.SystemConfiguration `json:"configuration"`
}
