package cmd

import "github.com/spf13/cobra"

// NewAclRoleCommand builds the "role" command group, which gathers the role management subcommands.
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
