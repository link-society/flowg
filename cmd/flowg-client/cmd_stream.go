package main

import (
	"context"

	"github.com/spf13/cobra"
)

func NewStreamCommand() *cobra.Command {
	type options struct {
		name string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "stream",
		Short: "Stream logs",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ctx = context.WithValue(ctx, StreamName, opts.name)
			cmd.SetContext(ctx)
		},
	}

	cmd.Flags().StringVar(
		&opts.name,
		"name",
		"default",
		"Name of the stream",
	)

	cmd.AddCommand(
		NewStreamWatchCommand(),
		NewStreamHistoryCommand(),
		NewStreamTailCommand(),
	)

	return cmd
}
