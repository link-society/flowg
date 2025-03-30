package main

import (
	"fmt"
	"os"

	"encoding/json"
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/client"
	"link-society.com/flowg/internal/utils/client/log"
	"link-society.com/flowg/internal/utils/client/sse"
)

func NewStreamWatchCommand() *cobra.Command {
	type options struct {
		name   string
		filter string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch logs in real-time",
		Run: func(cmd *cobra.Command, args []string) {
			url := fmt.Sprintf("/api/v1/streams/%s/logs/watch", opts.name)
			client := cmd.Context().Value(ApiClient).(*client.Client)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not prepare request: %v\n", err)
				exitCode = 1
				return
			}

			queryset := req.URL.Query()

			if opts.filter != "" {
				queryset.Set("filter", opts.filter)
			}

			req.URL.RawQuery = queryset.Encode()

			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not perform request: %v\n", err)
				exitCode = 1
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Fprintf(os.Stderr, "ERROR: Received non-200 response: %s\n", resp.Status)
				exitCode = 1
				return
			}

			printer := log.NewPrinter()
			stream := sse.NewEventStreamReader(resp.Body)
			for {
				event, err := stream.Next()
				if err == io.EOF {
					break
				}

				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not read event: %v\n", err)
					exitCode = 1
					return
				}

				switch event.Type {
				case "log":
					var log models.LogRecord

					if err := json.Unmarshal([]byte(event.Data), &log); err != nil {
						fmt.Fprintf(os.Stderr, "ERROR: Could not unmarshal log: %v\n", err)
						exitCode = 1
						return
					}

					if err := printer.Print(log); err != nil {
						fmt.Fprintf(os.Stderr, "ERROR: Could not print log: %v\n", err)
						exitCode = 1
						return
					}

				case "error":
					fmt.Fprintf(os.Stderr, "ERROR: %s\n", event.Data)
					exitCode = 1
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

	return cmd
}
