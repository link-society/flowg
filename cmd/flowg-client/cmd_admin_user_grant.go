package main

import (
	"fmt"
	"os"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/api"

	"link-society.com/flowg/internal/utils/client"
)

func NewAdminUserGrantCommand() *cobra.Command {
	type options struct {
		username string
		rolename string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "grant",
		Short: "Grant a role to a user",
		Run: func(cmd *cobra.Command, args []string) {
			client := cmd.Context().Value(ApiClient).(*client.Client)

			url := fmt.Sprintf("/api/v1/users/%s", opts.username)
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

			var data api.GetUserResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
				exitCode = 1
				return
			}

			found := false
			for _, role := range data.User.Roles {
				if role == opts.rolename {
					found = true
					break
				}
			}
			if !found {
				data.User.Roles = append(data.User.Roles, opts.rolename)

				body := struct {
					Roles []string `json:"roles"`
				}{
					Roles: data.User.Roles,
				}

				payload, err := json.Marshal(body)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not encode request body: %v\n", err)
					exitCode = 1
					return
				}

				req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
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
		&opts.username,
		"username",
		"",
		"Name of the user to grant the role to",
	)
	cmd.MarkFlagRequired("username")

	cmd.Flags().StringVar(
		&opts.rolename,
		"rolename",
		"",
		"Name of the role to grant to the user",
	)
	cmd.MarkFlagRequired("rolename")

	return cmd
}
