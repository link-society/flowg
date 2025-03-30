package main

import (
	"fmt"
	"os"

	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/api"

	"link-society.com/flowg/internal/utils/client"
)

func NewTransformerListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List transformers",
		Run: func(cmd *cobra.Command, args []string) {
			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := "/api/v1/transformers"
			req, err := http.NewRequest("GET", url, nil)
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

			var data api.ListTransformersResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
				exitCode = 1
				return
			}

			for _, name := range data.Transformers {
				fmt.Println(name)
			}
		},
	}
}
