package main

import (
	"log/slog"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/auth"
)

type adminDeleteRoleOpts struct {
	authDir string
	name    string
}

func NewAdminDeleteRoleCommand() *cobra.Command {
	opts := &adminDeleteRoleOpts{}

	cmd := &cobra.Command{
		Use:   "deleterole",
		Short: "Delete an existing role",
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

			err = authDb.DeleteRole(opts.name)
			if err != nil {
				slog.Error(
					"Failed to delete role",
					"channel", "main",
					"role", opts.name,
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
