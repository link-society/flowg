package main

import (
	"log/slog"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/auth"
)

type adminListRolesOpts struct {
	authDir string
}

func NewAdminListRolesCommand() *cobra.Command {
	opts := &adminListRolesOpts{}

	cmd := &cobra.Command{
		Use:   "listroles",
		Short: "List existing roles",
		Run: func(cmd *cobra.Command, args []string) {
			authDb, err := auth.NewDatabase(opts.authDir)
			if err != nil {
				slog.Error(
					"Failed to open auth database",
					"channel", "main",
					"path", opts.authDir,
					"error", err,
				)
				exitCode = 1
				return
			}
			defer func() {
				err := authDb.Close()
				if err != nil {
					slog.Error(
						"Failed to close auth database",
						"channel", "main",
						"path", opts.authDir,
						"error", err,
					)
					exitCode = 1
				}
			}()

			roleNames, err := authDb.ListRoles()
			if err != nil {
				slog.Error(
					"Failed to list roles",
					"channel", "main",
					"error", err,
				)
				exitCode = 1
				return
			}

			for _, roleName := range roleNames {
				role, err := authDb.GetRole(roleName)
				if err != nil {
					slog.Error(
						"Failed to get role",
						"channel", "main",
						"role", roleName,
						"error", err,
					)
					exitCode = 1
					return
				}

				slog.Info(
					"Role",
					"channel", "main",
					"name", role.Name,
					"scopes", role.Scopes,
				)
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.authDir,
		"auth-dir",
		"./data/auth",
		"Path to the log database directory",
	)
	cmd.MarkFlagDirname("auth-dir")

	return cmd
}
