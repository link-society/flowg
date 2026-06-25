package cmd

import (
	"github.com/spf13/cobra"
)

// NewStreamCommand builds the "stream" command group, which gathers the stream subcommands.
func NewStreamCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stream",
		Short: "Stream logs",
	}

	cmd.AddCommand(
		NewStreamListCommand(),
		NewStreamWatchCommand(),
		NewStreamHistoryCommand(),
		NewStreamTailCommand(),
		NewStreamSetCommand(),
		NewStreamIndexCommand(),
		NewStreamPurgeCommand(),
	)

	return cmd
}
