package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/logging"
)

var exitCode int = 0

func main() {
	var verbose bool

	rootCmd := &cobra.Command{
		Use:   "flowg",
		Short: "Low-Code log management solution",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logging.Setup(verbose)
		},
	}

	rootCmd.PersistentFlags().BoolVar(
		&verbose,
		"verbose",
		false,
		"Enable verbose logging",
	)

	rootCmd.AddCommand(
		NewServeCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
