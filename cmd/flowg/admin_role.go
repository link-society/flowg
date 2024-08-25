package main

import "github.com/spf13/cobra"

func NewAdminRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Role related admin commands (please run while the server is down)",
	}

	cmd.AddCommand(
		NewAdminRoleListCommand(),
		NewAdminRoleCreateCommand(),
		NewAdminRoleDeleteCommand(),
	)

	return cmd
}
