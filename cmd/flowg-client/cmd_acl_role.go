package main

import "github.com/spf13/cobra"

func NewAclRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Manage roles",
	}

	cmd.AddCommand(
		NewAclRoleListCommand(),
		NewAclRoleAddCommand(),
		NewAclRoleDeleteCommand(),
		NewAclRoleGrantCommand(),
		NewAclRoleRevokeCommand(),
	)

	return cmd
}
