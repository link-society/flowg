package cmd

import (
	"github.com/spf13/cobra"
)

// NewTokenCommand builds the "token" command group, which gathers the Personal Access Token subcommands.
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
