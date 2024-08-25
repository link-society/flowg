package main

import "github.com/spf13/cobra"

func NewAdminRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "role commands (please run while the server is down)",
	}

	cmd.AddCommand(
		NewAdminRoleListCommand(),
		NewAdminRoleCreateCommand(),
		NewAdminRoleDeleteCommand(),
	)

	return cmd
}
