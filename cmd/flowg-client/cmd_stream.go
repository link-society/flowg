package main

import (
	"github.com/spf13/cobra"
)

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
