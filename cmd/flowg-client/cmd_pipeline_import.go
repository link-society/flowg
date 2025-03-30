package main

import (
	"fmt"
	"os"

	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/client"
)

func NewPipelineImportCommand() *cobra.Command {
	type options struct {
		name string
		file string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			var reader io.Reader

			if opts.file != "" {
				file, err := os.Open(opts.file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not open file: %v\n", err)
					exitCode = 1
					return
				}
				defer file.Close()
				reader = file
			} else {
				reader = os.Stdin
			}

			body := struct {
				Flow models.FlowGraphV2
			}{}

			if err := json.NewDecoder(reader).Decode(&body.Flow); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode pipeline JSON: %v\n", err)
				exitCode = 1
				return
			}

			payload, err := json.Marshal(body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not encode request body: %v\n", err)
				exitCode = 1
				return
			}

			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := fmt.Sprintf("/api/v1/pipelines/%s", opts.name)
			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
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
		},
	}

	cmd.Flags().StringVar(
		&opts.name,
		"name",
		"default",
		"Name of the pipeline to import",
	)

	cmd.Flags().StringVar(
		&opts.file,
		"input",
		"",
		"Path to the pipeline file in JSON format (leave empty to read from stdin)",
	)
	cmd.MarkFlagFilename("file", "json")

	return cmd
}
