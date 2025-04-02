package main

import (
	"fmt"
	"os"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/api"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/client"
)

func NewAdminRoleGrantCommand() *cobra.Command {
	type options struct {
		name       string
		permission string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "grant",
		Short: "Grant a permission to a role",
		Run: func(cmd *cobra.Command, args []string) {
			newScope, err := models.ParseScope(opts.permission)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Invalid permission: %v\n", err)
				exitCode = 1
				return
			}

			client := cmd.Context().Value(ApiClient).(*client.Client)

			url := fmt.Sprintf("/api/v1/roles/%s", opts.name)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not prepare request: %v\n", err)
				exitCode = 1
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not send request: %v\n", err)
				exitCode = 1
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Fprintf(os.Stderr, "ERROR: Received non-200 response: %s\n", resp.Status)
				exitCode = 1
				return
			}

			var data api.GetRoleResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
				exitCode = 1
				return
			}

			found := false
			for _, scope := range data.Role.Scopes {
				if scope == newScope {
					found = true
					break
				}
			}
			if !found {
				data.Role.Scopes = append(data.Role.Scopes, newScope)
				scopes := make([]string, len(data.Role.Scopes))
				for i, scope := range data.Role.Scopes {
					scopes[i] = string(scope)
				}

				body := struct {
					Scopes []string `json:"scopes"`
				}{
					Scopes: scopes,
				}

				payload, err := json.Marshal(body)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not encode request body: %v\n", err)
					exitCode = 1
					return
				}

				req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not prepare request: %v\n", err)
					exitCode = 1
					return
				}

				resp, err = client.Do(req)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not send request: %v\n", err)
					exitCode = 1
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					fmt.Fprintf(os.Stderr, "ERROR: Received non-200 response: %s\n", resp.Status)
					exitCode = 1
					return
				}
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.name,
		"name",
		"",
		"Name of the role to grant the permission to",
	)
	cmd.MarkFlagRequired("name")

	cmd.Flags().StringVar(
		&opts.permission,
		"permission",
		"",
		"Permission to grant to the role",
	)
	cmd.MarkFlagRequired("permission")

	return cmd
}
