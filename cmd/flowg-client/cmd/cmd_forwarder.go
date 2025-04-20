package cmd

import "github.com/spf13/cobra"

func NewForwarderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "forwarder",
		Short: "Access forwarders",
	}

	cmd.AddCommand(
		NewForwarderListCommand(),
		NewForwarderExportCommand(),
		NewForwarderImportCommand(),
		NewForwarderDeleteCommand(),
	)

	return cmd
}
