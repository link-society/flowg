package cmd

import "github.com/spf13/cobra"

// NewTransformerCommand builds the "transformer" command group, which gathers the transformer subcommands.
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
