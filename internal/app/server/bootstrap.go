package server

import (
	"context"
	"fmt"
	"log/slog"

	"link-society.com/flowg/internal/storage"

	"link-society.com/flowg/internal/storage/bootstrap"
)

type bootstrapHandler struct {
	logger *slog.Logger

	authStorage   storage.AuthStorage
	configStorage storage.ConfigStorage

	initialSyslogAllowedOrigins []string

	initialUser     string
	initialPassword string

	resetUser     string
	resetPassword string
}

func (h *bootstrapHandler) Run(ctx context.Context) error {
	err := bootstrap.DefaultSystemConfig(ctx, h.configStorage, bootstrap.BootstrapSystemOptions{
		InitialSyslogAllowedOrigins: h.initialSyslogAllowedOrigins,
	})
	if err != nil {
		return fmt.Errorf("failed to bootstrap default system config: %w", err)
	}

	err = bootstrap.DefaultRolesAndUsers(ctx, h.authStorage, bootstrap.BootstrapAuthOptions{
		InitialUser:     h.initialUser,
		InitialPassword: h.initialPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to bootstrap default roles and users: %w", err)
	}

	err = bootstrap.DefaultPipeline(ctx, h.configStorage)
	if err != nil {
		return fmt.Errorf("failed to bootstrap default pipeline: %w", err)
	}

	err = bootstrap.ResetUser(ctx, h.authStorage, bootstrap.ResetUserOptions{
		User:     h.resetUser,
		Password: h.resetPassword,
	})

	return nil
}
