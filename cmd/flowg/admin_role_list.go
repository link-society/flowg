package main

import (
	"context"
	"time"

	"os"

	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/storage/auth"
)

type adminRoleListOpts struct {
	authDir string
}

func NewAdminRoleListCommand() *cobra.Command {
	opts := &adminRoleListOpts{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List existing roles",
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

			roles, err := authStorage.ListRoles(context.Background())
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR: Failed to list roles:", err)
				exitCode = 1
				return
			}

			if len(roles) == 0 {
				fmt.Println("No roles found")
			} else {
				writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

				fmt.Fprintln(writer, "Name\tScopes")

				for _, role := range roles {
					fmt.Fprintf(writer, "%s\t%v\n", role.Name, role.Scopes)
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
