package main

import "github.com/spf13/cobra"

func NewAdminTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Personal Access Token commands (please run while the server is down)",
	}

	cmd.AddCommand(
		NewAdminTokenCreateCommand(),
	)

	return cmd
}
