package main

import (
	"fmt"

	"context"
	"time"

	"os"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/storage/auth"
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
			authStorage := auth.NewStorage(auth.OptDirectory(opts.authDir))
			authStorage.Start()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := authStorage.WaitReady(ctx)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to open auth database:", err)
				exitCode = 1
				return
			}

			defer func() {
				authStorage.Stop()

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				err := authStorage.Join(ctx)
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to close auth database:", err)
					exitCode = 1
				}
			}()

			err = authStorage.DeleteUser(context.Background(), opts.name)
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
