package cmd

import "github.com/spf13/cobra"

// NewAdminCommand builds the "admin" command group, which gathers the admin subcommands.
func NewAdminCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Admin commands",
	}

	cmd.AddCommand(
		NewAdminBackupCommand(),
		NewAdminRestoreCommand(),
	)

	return cmd
}
