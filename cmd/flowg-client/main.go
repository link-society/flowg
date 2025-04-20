package main

import (
	"os"

	"link-society.com/flowg/cmd/flowg-client/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()

	if err := rootCmd.Execute(); err != nil {
		cmd.ExitCode = 1
	}

	os.Exit(cmd.ExitCode)
}
