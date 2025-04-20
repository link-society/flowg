package cmd

import (
	"bytes"
	"fmt"
	"os"

	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/utils/client"
)

func NewAclRoleAddCommand() *cobra.Command {
	type options struct {
		name string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add new role",
		Run: func(cmd *cobra.Command, args []string) {
			body := struct {
				Scopes []string `json:"scopes"`
			}{
				Scopes: []string{},
			}

			payload, err := json.Marshal(body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to marshal request body: %v\n", err)
				ExitCode = 1
				return
			}

			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := fmt.Sprintf("/api/v1/roles/%s", opts.name)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to prepare request: %v\n", err)
				ExitCode = 1
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to send request: %v\n", err)
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
		&opts.name,
		"name",
		"",
		"Name of the new role",
	)
	cmd.MarkFlagRequired("name")

	return cmd
}
