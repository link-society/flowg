package main

import (
	"os"

	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/data/auth"
)

type adminUserListOpts struct {
	authDir string
}

func NewAdminUserListCommand() *cobra.Command {
	opts := &adminUserListOpts{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List existing users",
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
			users, err := userSys.ListUsers()
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to list users:", err)
				exitCode = 1
				return
			}

			if len(users) == 0 {
				fmt.Println("No users found")
			} else {
				writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

				fmt.Fprintln(writer, "Name\tRoles")

				for _, user := range users {
					fmt.Fprintf(writer, "%s\t%v\n", user.Name, user.Roles)
				}

				writer.Flush()
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.authDir,
		"auth-dir",
		defaultAuthDir,
		"Path to the log database directory",
	)
	cmd.MarkFlagDirname("auth-dir")

	return cmd
}
