package main

import (
	"fmt"
	"os"

	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/utils/client"
)

func NewAdminRoleDeleteCommand() *cobra.Command {
	type options struct {
		name string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove role",
		Run: func(cmd *cobra.Command, args []string) {
			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := fmt.Sprintf("/api/v1/roles/%s", opts.name)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to prepare request: %v\n", err)
				exitCode = 1
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to send request: %v\n", err)
				exitCode = 1
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Fprintf(os.Stderr, "ERROR: Received non-200 response: %s\n", resp.Status)
				exitCode = 1
				return
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.name,
		"name",
		"",
		"Name of the role",
	)
	cmd.MarkFlagRequired("name")

	return cmd
}
