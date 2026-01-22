package cmd

import (
	"github.com/spf13/cobra"
)

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
