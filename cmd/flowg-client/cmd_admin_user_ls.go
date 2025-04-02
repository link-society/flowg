package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/api"

	"link-society.com/flowg/internal/utils/client"
)

func NewAdminUserListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List all users",
		Run: func(cmd *cobra.Command, args []string) {
			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := "/api/v1/users"
			req, err := http.NewRequest(http.MethodGet, url, nil)
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

			var data api.ListUsersResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to decode response: %v\n", err)
				exitCode = 1
				return
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
			fmt.Fprintf(w, "USERNAME\tROLES\n")

			for _, user := range data.Users {
				roles := strings.Join(user.Roles, ",")
				fmt.Fprintf(w, "%s\t%s\n", user.Name, roles)
			}

			if err := w.Flush(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to flush output: %v\n", err)
				exitCode = 1
				return
			}
		},
	}

	return cmd
}
