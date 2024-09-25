package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/data/auth"
)

type adminUserDeleteOpts struct {
	authDir string
	name    string
}

func NewAdminUserDeleteCommand() *cobra.Command {
	opts := &adminUserDeleteOpts{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an existing user",
		Run: func(cmd *cobra.Command, args []string) {
			authDb := auth.NewDatabase(
				auth.DefaultDatabaseOpts().WithDir(opts.authDir),
			)
			err := authDb.Open()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open auth database:", err)
				exitCode = 1
				return
			}
			defer func() {
				err := authDb.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close auth database:", err)
					exitCode = 1
				}
			}()

			userSys := auth.NewUserSystem(authDb)
			err = userSys.DeleteUser(opts.name)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to delete user:", err)
				exitCode = 1
				return
			}

			fmt.Println("User deleted")
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
		"Name of the user",
	)
	cmd.MarkFlagRequired("name")

	return cmd
}
