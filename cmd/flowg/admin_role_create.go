package main

import (
	"os"

	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/auth"
)

type adminRoleCreateOpts struct {
	authDir string
	name    string
}

func NewAdminRoleCreateCommand() *cobra.Command {
	opts := &adminRoleCreateOpts{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new role",
		Run: func(cmd *cobra.Command, args []string) {
			role := auth.Role{
				Name:   opts.name,
				Scopes: make([]auth.Scope, len(args)),
			}

			for i, scopeName := range args {
				scope, err := auth.ParseScope(scopeName)
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR: Failed to parse scope:", err)
					exitCode = 1
					return
				}
				role.Scopes[i] = scope
			}

			authDb, err := auth.NewDatabase(opts.authDir)
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

			err = authDb.SaveRole(role)
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
