package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/storage/auth"
)

type adminRoleDeleteOpts struct {
	authDir string
	name    string
}

func NewAdminRoleDeleteCommand() *cobra.Command {
	opts := &adminRoleDeleteOpts{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an existing role",
		Run: func(cmd *cobra.Command, args []string) {
			authStorage := auth.NewStorage(
				auth.OptDirectory(opts.authDir),
			)
			authStorage.Start()
			err := authStorage.WaitStarted()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open auth database:", err)
				exitCode = 1
				return
			}
			defer func() {
				authStorage.Stop()
				err := authStorage.WaitStopped()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close auth database:", err)
					exitCode = 1
				}
			}()

			ctx := context.Background()
			err = authStorage.DeleteRole(ctx, opts.name)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to delete role:", err)
				exitCode = 1
				return
			}

			fmt.Println("Role deleted")
		},
	}

	cmd.Flags().StringVar(
		&opts.authDir,
		"auth-dir",
		defaultAuthDir,
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
