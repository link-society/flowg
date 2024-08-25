package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/logging"
)

var exitCode int = 0

func main() {
	rootCmd := &cobra.Command{
		Use:   "flowg",
		Short: "Low-Code log management solution",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logging.Discard()
		},
	}

	rootCmd.AddCommand(
		NewAdminCommand(),
		NewServeCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
