package cmd

import "github.com/spf13/cobra"

func NewAclCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "Access Control commands",
	}

	cmd.AddCommand(
		NewAclUserCommand(),
		NewAclRoleCommand(),
	)

	return cmd
}
