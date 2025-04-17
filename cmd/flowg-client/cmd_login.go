package main

import (
	"bytes"
	"fmt"
	"strings"

	"encoding/json"
	"net/http"

	"io"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"link-society.com/flowg/api"
	"link-society.com/flowg/internal/utils/client"
)

func NewLoginCommand() *cobra.Command {
	type options struct {
		username string
		password string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Get a temporary JSON Web Token for authentication",
		Run: func(cmd *cobra.Command, args []string) {
			switch opts.password {
			case "":
				fmt.Print("Password: ")
				passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
				fmt.Println()
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to read password from prompt: %v\n", err)
					exitCode = 1
					return
				}

				opts.password = strings.TrimSpace(string(passwordBytes))

			case "-":
				passwordBytes, err := io.ReadAll(os.Stdin)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to read password from stdin: %v\n", err)
					exitCode = 1
					return
				}

				opts.password = strings.TrimSpace(string(passwordBytes))
			}

			body := struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}{
				Username: opts.username,
				Password: opts.password,
			}
			payload, err := json.Marshal(body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to marshal request body: %v\n", err)
				exitCode = 1
				return
			}

			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := "/api/v1/auth/login"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not prepare request: %v\n", err)
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

			var data api.LoginResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to decode response: %v\n", err)
				exitCode = 1
				return
			}

			fmt.Println("Token:", data.Token)
		},
	}

	cmd.Flags().StringVar(
		&opts.username,
		"username",
		"",
		"Username to authenticate as",
	)
	cmd.MarkFlagRequired("username")

	cmd.Flags().StringVar(
		&opts.password,
		"password",
		"",
		"Password to authenticate with (if empty, will be prompted, set to - to read from stdin)",
	)

	return cmd
}
