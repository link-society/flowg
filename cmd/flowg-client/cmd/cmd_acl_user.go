package cmd

import "github.com/spf13/cobra"

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
