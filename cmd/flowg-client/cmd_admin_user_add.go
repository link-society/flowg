package main

import (
	"bytes"
	"fmt"
	"os"

	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/utils/client"
)

func NewAdminUserAddCommand() *cobra.Command {
	type options struct {
		username string
		password string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add new user",
		Run: func(cmd *cobra.Command, args []string) {
			body := struct {
				Roles    []string `json:"roles"`
				Password string   `json:"password"`
			}{
				Roles:    []string{},
				Password: opts.password,
			}

			payload, err := json.Marshal(body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to marshal request body: %v\n", err)
				exitCode = 1
				return
			}

			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := fmt.Sprintf("/api/v1/users/%s", opts.username)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
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
		&opts.username,
		"username",
		"",
		"Name of the new user",
	)
	cmd.MarkFlagRequired("username")

	cmd.Flags().StringVar(
		&opts.password,
		"password",
		"",
		"Password for the new user",
	)
	cmd.MarkFlagRequired("password")

	return cmd
}
