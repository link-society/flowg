package server

import (
	"context"
	"fmt"
	"log/slog"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"

	"link-society.com/flowg/internal/app/bootstrap"
)

type bootstrapHandler struct {
	logger *slog.Logger

	authStorage   auth.Storage
	configStorage config.Storage

	initialUser     string
	initialPassword string

	resetUser     string
	resetPassword string
}

func (h *bootstrapHandler) Run(ctx context.Context) error {
	err := bootstrap.DefaultRolesAndUsers(ctx, h.authStorage, bootstrap.BootstrapOptions{
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
