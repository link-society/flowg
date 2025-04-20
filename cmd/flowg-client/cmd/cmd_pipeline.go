package cmd

import "github.com/spf13/cobra"

func NewPipelineCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Access pipelines",
	}

	cmd.AddCommand(
		NewPipelineListCommand(),
		NewPipelineExportCommand(),
		NewPipelineImportCommand(),
		NewPipelineDeleteCommand(),
	)

	return cmd
}
