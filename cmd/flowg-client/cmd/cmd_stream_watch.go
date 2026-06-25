package cmd

import (
	"fmt"
	"os"

	"encoding/json"
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/cmd/flowg-client/utils"
	"link-society.com/flowg/cmd/flowg-client/utils/sse"
)

// NewStreamWatchCommand builds the "watch" command, which watches logs in real-time.
func NewStreamWatchCommand() *cobra.Command {
	type options struct {
		name     string
		filter   string
		indexing utils.IndexMap
	}

	opts := &options{
		indexing: make(utils.IndexMap),
	}

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch logs in real-time",
		Run: func(cmd *cobra.Command, args []string) {
			url := fmt.Sprintf("/api/v1/streams/%s/logs/watch", opts.name)
			client := cmd.Context().Value(ApiClient).(*utils.Client)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not prepare request: %v\n", err)
				ExitCode = 1
				return
			}

			queryset := req.URL.Query()

			if opts.filter != "" {
				queryset.Set("filter", opts.filter)
			}

			if len(opts.indexing) > 0 {
				payload, err := json.Marshal(opts.indexing)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not encode indexing parameters: %v\n", err)
					ExitCode = 1
					return
				}

				queryset.Set("indexing", string(payload))
			}

			req.URL.RawQuery = queryset.Encode()

			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not perform request: %v\n", err)
				ExitCode = 1
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Fprintf(os.Stderr, "ERROR: Received non-200 response: %s\n", resp.Status)
				ExitCode = 1
				return
			}

			printer := utils.NewPrinter()
			stream := sse.NewEventStreamReader(resp.Body)
			for {
				event, err := stream.Next()
				if err == io.EOF {
					break
				}

				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not read event: %v\n", err)
					ExitCode = 1
					return
				}

				switch event.Type {
				case "log":
					var log models.LogRecord

					if err := json.Unmarshal([]byte(event.Data), &log); err != nil {
						fmt.Fprintf(os.Stderr, "ERROR: Could not unmarshal log: %v\n", err)
						ExitCode = 1
						return
					}

					if err := printer.Print(log); err != nil {
						fmt.Fprintf(os.Stderr, "ERROR: Could not print log: %v\n", err)
						ExitCode = 1
						return
					}

				case "error":
					fmt.Fprintf(os.Stderr, "ERROR: %s\n", event.Data)
					ExitCode = 1
					return
				}
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.name,
		"name",
		"default",
		"Name of the stream",
	)

	cmd.Flags().StringVar(
		&opts.filter,
		"filter",
		"",
		"Filter logs",
	)

	cmd.Flags().Var(
		&opts.indexing,
		"index",
		"Indexing key-value pairs to filter logs (can be specified multiple times)",
	)

	return cmd
}
