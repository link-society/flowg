package cmd

import "github.com/spf13/cobra"

// NewAclCommand builds the "acl" command group, which gathers the Access Control subcommands.
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
