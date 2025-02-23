package main

import (
	"context"
	"time"

	"os"

	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
)

type adminRoleCreateOpts struct {
	authDir string
	name    string
}

func NewAdminRoleCreateCommand() *cobra.Command {
	opts := &adminRoleCreateOpts{}

	cmd := &cobra.Command{
		Use:   "create [flags] [...scopes]",
		Short: "Create a new role",
		Run: func(cmd *cobra.Command, args []string) {
			role := models.Role{
				Name:   opts.name,
				Scopes: make([]models.Scope, len(args)),
			}

			for i, scopeName := range args {
				scope, err := models.ParseScope(scopeName)
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to parse scope:", err)
					exitCode = 1
					return
				}
				role.Scopes[i] = scope
			}

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

			err = authStorage.SaveRole(context.Background(), role)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to save role:", err)
				exitCode = 1
				return
			}

			writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

			fmt.Fprintln(writer, "Name\tScopes")
			fmt.Fprintf(writer, "%s\t%s\n", role.Name, role.Scopes)

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
		"Name of the role",
	)
	cmd.MarkFlagRequired("name")

	return cmd
}
