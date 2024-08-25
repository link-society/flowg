package main

import "github.com/spf13/cobra"

func NewAdminCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Admin commands (please run while the server is down)",
	}

	cmd.AddCommand(
		NewAdminRoleCommand(),
		NewAdminUserCommand(),
	)

	return cmd
}
