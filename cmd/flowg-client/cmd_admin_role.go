package main

import "github.com/spf13/cobra"

func NewAdminRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Manage roles",
	}

	cmd.AddCommand(
		NewAdminRoleListCommand(),
		NewAdminRoleAddCommand(),
		NewAdminRoleDeleteCommand(),
		NewAdminRoleGrantCommand(),
		NewAdminRoleRevokeCommand(),
	)

	return cmd
}
