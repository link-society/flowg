package cmd

import "github.com/spf13/cobra"

func NewTransformerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transformer",
		Short: "Access transformers",
	}

	cmd.AddCommand(
		NewTransformerListCommand(),
		NewTransformerCatCommand(),
		NewTransformerEditCommand(),
		NewTransformerDeleteCommand(),
	)

	return cmd
}
