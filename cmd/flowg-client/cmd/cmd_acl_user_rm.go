package cmd

import (
	"fmt"
	"os"

	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/cmd/flowg-client/utils"
)

// NewAclUserDeleteCommand builds the "rm" command, which removes a user.
func NewAclUserDeleteCommand() *cobra.Command {
	type options struct {
		username string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove user",
		Run: func(cmd *cobra.Command, args []string) {
			client := cmd.Context().Value(ApiClient).(*utils.Client)
			url := fmt.Sprintf("/api/v1/users/%s", opts.username)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
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
		&opts.username,
		"username",
		"",
		"Name of the user",
	)
	cmd.MarkFlagRequired("username")

	return cmd
}
