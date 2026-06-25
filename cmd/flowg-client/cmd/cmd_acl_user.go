package cmd

import "github.com/spf13/cobra"

// NewAclUserCommand builds the "user" command group, which gathers the user management subcommands.
func NewAclUserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}

	cmd.AddCommand(
		NewAclUserListCommand(),
		NewAclUserAddCommand(),
		NewAclUserDeleteCommand(),
		NewAclUserGrantCommand(),
		NewAclUserRevokeCommand(),
	)

	return cmd
}
