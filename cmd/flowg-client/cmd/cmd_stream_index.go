package cmd

import (
	"fmt"
	"os"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/api"

	"slices"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/client"
)

func NewStreamIndexCommand() *cobra.Command {
	type options struct {
		name   string
		field  string
		remove bool
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "index",
		Short: "Create or remove indexed fields for a stream",
		Run: func(cmd *cobra.Command, args []string) {
			url := fmt.Sprintf("/api/v1/streams/%s", opts.name)
			client := cmd.Context().Value(ApiClient).(*client.Client)
			req, err := http.NewRequest(http.MethodGet, url, nil)
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
				fmt.Fprintf(os.Stderr, "ERROR: Failed to get stream properties: %s\n", resp.Status)
				ExitCode = 1
				return
			}

			var data api.GetStreamResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
				ExitCode = 1
				return
			}

			streamConfig := data.Config

			if opts.remove {
				for i, field := range streamConfig.IndexedFields {
					if field == opts.field {
						streamConfig.IndexedFields = append(
							streamConfig.IndexedFields[:i],
							streamConfig.IndexedFields[i+1:]...,
						)
						break
					}
				}
			} else {
				if !slices.Contains(streamConfig.IndexedFields, opts.field) {
					streamConfig.IndexedFields = append(streamConfig.IndexedFields, opts.field)
				}
			}

			body := struct {
				Config models.StreamConfig `json:"config"`
			}{
				Config: streamConfig,
			}
			payload, err := json.Marshal(body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not encode request body: %v\n", err)
				ExitCode = 1
				return
			}

			req, err = http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not prepare request: %v\n", err)
				ExitCode = 1
				return
			}

			resp, err = client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not send request: %v\n", err)
				ExitCode = 1
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to set stream properties: %s\n", resp.Status)
				ExitCode = 1
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

	cmd.Flags().StringVar(
		&opts.field,
		"field",
		"",
		"Name of the field to (un)index",
	)
	cmd.MarkFlagRequired("field")

	cmd.Flags().BoolVar(
		&opts.remove,
		"rm",
		false,
		"If set, removes the field from the index",
	)

	return cmd
}
