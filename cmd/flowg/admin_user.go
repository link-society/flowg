package main

import "github.com/spf13/cobra"

func NewAdminUserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "User commands (please run while the server is down)",
	}

	cmd.AddCommand(
		NewAdminUserListCommand(),
		NewAdminUserCreateCommand(),
		NewAdminUserDeleteCommand(),
	)

	return cmd
}
