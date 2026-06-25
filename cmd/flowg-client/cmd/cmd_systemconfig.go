package cmd

import (
	"github.com/spf13/cobra"
)

// NewSystemConfigCommand builds the "system-config" command group, which gathers the system configuration subcommands.
func NewSystemConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system-config",
		Short: "System configuration commands",
	}

	cmd.AddCommand(
		NewSystemConfigShowCommand(),
		NewSystemConfigUpdateCommand(),
	)

	return cmd
}
