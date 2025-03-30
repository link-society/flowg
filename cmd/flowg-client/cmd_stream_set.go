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

func NewStreamSetCommand() *cobra.Command {
	type options struct {
		name          string
		retentionTime int64
		retentionSize int64
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set stream properties",
		Run: func(cmd *cobra.Command, args []string) {
			url := fmt.Sprintf("/api/v1/streams/%s", opts.name)
			client := cmd.Context().Value(ApiClient).(*client.Client)
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
				fmt.Fprintf(os.Stderr, "ERROR: Failed to get stream properties: %s\n", resp.Status)
				exitCode = 1
				return
			}

			var data api.GetStreamResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
				exitCode = 1
				return
			}

			streamConfig := data.Config

			if cmd.Flags().Changed("ttl") {
				streamConfig.RetentionTime = opts.retentionTime
			}

			if cmd.Flags().Changed("max-size") {
				streamConfig.RetentionSize = opts.retentionSize
			}

			body := struct {
				Config models.StreamConfig `json:"config"`
			}{
				Config: streamConfig,
			}
			payload, err := json.Marshal(body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not encode request body: %v\n", err)
				exitCode = 1
				return
			}

			req, err = http.NewRequest("PUT", url, bytes.NewBuffer(payload))
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
				fmt.Fprintf(os.Stderr, "ERROR: Failed to set stream properties: %s\n", resp.Status)
				exitCode = 1
				return
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.name,
		"name",
		"default",
		"Name of the stream to set properties for",
	)

	cmd.Flags().Int64Var(
		&opts.retentionTime,
		"ttl",
		0,
		"Retention time in seconds (0 to disable)",
	)

	cmd.Flags().Int64Var(
		&opts.retentionSize,
		"max-size",
		0,
		"Max size in bytes (0 to disable)",
	)

	return cmd
}
