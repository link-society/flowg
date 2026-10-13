package server

import (
	"context"
	"fmt"
	"log/slog"

	"go.uber.org/fx"

	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/storage/bootstrap"
)

// bootstrapOptions configures the bootstrap handler: the syslog origins allowed
// on a fresh instance, the initial admin credentials, whether to seed the demo
// role and user, and the optional credential reset.
type bootstrapOptions struct {
	InitialSyslogAllowedOrigins []string

	InitialUser     string
	InitialPassword string

	DemoMode bool

	ResetUser     string
	ResetPassword string
}

// bootstrapHandler seeds a fresh (or upgraded) instance on startup: it applies
// the default system configuration, roles and users, and pipeline, and performs
// the optional admin-credential reset requested through the CLI.
type bootstrapHandler struct {
	logger *slog.Logger

	authStorage   storage.AuthStorage
	configStorage storage.ConfigStorage

	opts bootstrapOptions
}

// newBootstrapHandler builds the bootstrap handler configured by the
// [bootstrapOptions] provided to the graph, logging through the root logger,
// and runs it from the OnStart hook.
// It depends on the log storage too, so that every storage is started before
// the handler runs.
func newBootstrapHandler(
	lc fx.Lifecycle,
	opts bootstrapOptions,
	rootLogger *slog.Logger,
	authStorage storage.AuthStorage,
	configStorage storage.ConfigStorage,
	_ storage.LogStorage,
) *bootstrapHandler {
	h := &bootstrapHandler{
		logger:        rootLogger.With(slog.String("channel", "bootstrap")),
		authStorage:   authStorage,
		configStorage: configStorage,
		opts:          opts,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return h.Run(ctx)
		},
	})

	return h
}

// Run applies all bootstrap steps in order: default roles and users,
// default system configuration, default pipeline, and the optional
// user-password reset. It is invoked from the module's OnStart lifecycle hook.
func (h *bootstrapHandler) Run(ctx context.Context) error {
	err := bootstrap.DefaultRolesAndUsers(ctx, h.authStorage, bootstrap.BootstrapAuthOptions{
		InitialUser:     h.opts.InitialUser,
		InitialPassword: h.opts.InitialPassword,
		DemoMode:        h.opts.DemoMode,
	})
	if err != nil {
		return fmt.Errorf("failed to bootstrap default roles and users: %w", err)
	}

	err = bootstrap.DefaultSystemConfig(ctx, h.configStorage, bootstrap.BootstrapSystemOptions{
		InitialSyslogAllowedOrigins: h.opts.InitialSyslogAllowedOrigins,
	})
	if err != nil {
		return fmt.Errorf("failed to bootstrap default system config: %w", err)
	}

	err = bootstrap.DefaultPipeline(ctx, h.configStorage)
	if err != nil {
		return fmt.Errorf("failed to bootstrap default pipeline: %w", err)
	}

	err = bootstrap.ResetUser(ctx, h.logger, h.authStorage, bootstrap.ResetUserOptions{
		User:     h.opts.ResetUser,
		Password: h.opts.ResetPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to reset user credentials: %w", err)
	}

	return nil
}
