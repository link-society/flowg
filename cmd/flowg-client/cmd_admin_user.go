package main

import "github.com/spf13/cobra"

func NewAdminUserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}

	cmd.AddCommand(
		NewAdminUserListCommand(),
		NewAdminUserAddCommand(),
		NewAdminUserDeleteCommand(),
		NewAdminUserGrantCommand(),
		NewAdminUserRevokeCommand(),
	)

	return cmd
}
