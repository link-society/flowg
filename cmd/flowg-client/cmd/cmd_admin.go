package cmd

import "github.com/spf13/cobra"

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
