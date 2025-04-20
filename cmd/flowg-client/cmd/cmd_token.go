package cmd

import (
	"github.com/spf13/cobra"
)

func NewTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage Personal Access Tokens",
	}

	cmd.AddCommand(
		NewTokenListCommand(),
		NewTokenCreateCommand(),
		NewTokenRevokeCommand(),
	)

	return cmd
}
