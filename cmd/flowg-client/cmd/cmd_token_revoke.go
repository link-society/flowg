package cmd

import (
	"fmt"
	"os"

	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/utils/client"
)

func NewTokenRevokeCommand() *cobra.Command {
	type options struct {
		uuid string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a Personal Access Token",
		Run: func(cmd *cobra.Command, args []string) {
			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := fmt.Sprintf("/api/v1/tokens/%s", opts.uuid)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not prepare request: %v\n", err)
				ExitCode = 1
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not send request: %v\n", err)
				ExitCode = 1
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Fprintf(os.Stderr, "ERROR: Received non-200 response: %s\n", resp.Status)
				ExitCode = 1
				return
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.uuid,
		"uuid",
		"",
		"UUID of the Personal Access Token to delete",
	)
	cmd.MarkFlagRequired("uuid")

	return cmd
}
