package main

import (
	"context"
	"os"

	"fmt"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/storage/auth"
)

type adminTokenCreateOpts struct {
	authDir string
	user    string
}

func NewAdminTokenCreateCommand() *cobra.Command {
	opts := &adminTokenCreateOpts{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Personal Access Token",
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
			user, err := authStorage.FetchUser(ctx, opts.user)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to get user:", err)
				exitCode = 1
				return
			}

			if user == nil {
				fmt.Fprintln(os.Stderr, "ERROR: User not found")
				exitCode = 1
				return
			}

			token, _, err := authStorage.CreateToken(ctx, user.Name)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to generate token:", err)
				exitCode = 1
				return
			}

			fmt.Println(token)
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
		&opts.user,
		"user",
		"",
		"Name of the user",
	)
	cmd.MarkFlagRequired("user")

	return cmd
}
