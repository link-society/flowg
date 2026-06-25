package bootstrap

import (
	"context"
	"fmt"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
)

type BootstrapSystemOptions struct {
	InitialSyslogAllowedOrigins []string
}

func DefaultSystemConfig(ctx context.Context, configStorage storage.ConfigStorage, opts BootstrapSystemOptions) error {
	hasConfig, err := configStorage.HasSystemConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if system config exists: %w", err)
	}

	if !hasConfig {
		defaultConfig := &models.SystemConfiguration{
			SyslogAllowedOrigins: opts.InitialSyslogAllowedOrigins,
		}

		if err := configStorage.WriteSystemConfig(ctx, defaultConfig); err != nil {
			return fmt.Errorf("failed to write default system config: %w", err)
		}
	}

	return nil
}
