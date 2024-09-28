package main

import (
	"context"
	"os"

	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
)

type adminUserCreateOpts struct {
	authDir  string
	name     string
	password string
}

func NewAdminUserCreateCommand() *cobra.Command {
	opts := &adminUserCreateOpts{}

	cmd := &cobra.Command{
		Use:   "create [flags] [...roles]",
		Short: "Create a new user",
		Run: func(cmd *cobra.Command, args []string) {
			user := models.User{
				Name:  opts.name,
				Roles: make([]string, len(args)),
			}

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

			copy(user.Roles, args)

			ctx := context.Background()
			err = authStorage.SaveUser(ctx, user, opts.password)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to save user:", err)
				exitCode = 1
				return
			}

			writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

			fmt.Fprintln(writer, "Name\tRoles")
			fmt.Fprintf(writer, "%s\t%v\n", user.Name, user.Roles)

			writer.Flush()
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

	cmd.Flags().StringVar(
		&opts.password,
		"password",
		"",
		"Password of the user",
	)
	cmd.MarkFlagRequired("password")

	return cmd
}
