package main

import (
	"log/slog"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/auth"
)

type adminCreateRoleOpts struct {
	authDir string
	name    string
}

func NewAdminCreateRoleCommand() *cobra.Command {
	opts := &adminCreateRoleOpts{}

	cmd := &cobra.Command{
		Use:   "createrole",
		Short: "Create a new role",
		Run: func(cmd *cobra.Command, args []string) {
			role := auth.Role{
				Name:   opts.name,
				Scopes: make([]auth.Scope, len(args)),
			}

			for i, scopeName := range args {
				scope, err := auth.ParseScope(scopeName)
				if err != nil {
					slog.Error(
						"Failed to parse scope",
						"channel", "main",
						"scope", scopeName,
						"error", err.Error(),
					)
					exitCode = 1
					return
				}
				role.Scopes[i] = scope
			}

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

			err = authDb.SaveRole(role)
			if err != nil {
				slog.Error(
					"Failed to save role",
					"channel", "main",
					"role", role.Name,
					"error", err,
				)
				exitCode = 1
				return
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

	cmd.Flags().StringVar(
		&opts.name,
		"name",
		"",
		"Name of the role",
	)
	cmd.MarkFlagRequired("name")

	return cmd
}
