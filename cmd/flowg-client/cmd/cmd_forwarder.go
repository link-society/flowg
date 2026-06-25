package cmd

import "github.com/spf13/cobra"

// NewForwarderCommand builds the "forwarder" command group, which gathers the forwarder subcommands.
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
