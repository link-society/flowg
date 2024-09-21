package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/app"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(*cobra.Command, []string) {
			fmt.Println(app.FLOWG_VERSION)
		},
	}
}
