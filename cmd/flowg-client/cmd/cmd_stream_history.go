package cmd

import (
	"fmt"
	"os"
	"time"

	"encoding/json"
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/api"

	"link-society.com/flowg/internal/utils/client"
	"link-society.com/flowg/internal/utils/client/flags"
	"link-society.com/flowg/internal/utils/client/log"
)

func NewStreamHistoryCommand() *cobra.Command {
	type options struct {
		name     string
		filter   string
		from     string
		to       string
		indexing flags.IndexMap
	}

	opts := &options{
		indexing: make(flags.IndexMap),
	}

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Fetch logs using a time window",
		Run: func(cmd *cobra.Command, args []string) {
			url := fmt.Sprintf("/api/v1/streams/%s/logs", opts.name)
			client := cmd.Context().Value(ApiClient).(*client.Client)
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

			queryset.Set("from", opts.from)
			queryset.Set("to", opts.to)

			req.URL.RawQuery = queryset.Encode()

			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not send request: %v\n", err)
				ExitCode = 1
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Fprintf(os.Stderr, "ERROR: Unexpected status code: %d\n", resp.StatusCode)
				io.Copy(os.Stderr, resp.Body)
				ExitCode = 1
				return
			}

			printer := log.NewPrinter()

			var data api.QueryStreamResponse

			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
				ExitCode = 1
				return
			}

			for _, log := range data.Records {
				if err := printer.Print(log); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not print log: %v\n", err)
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

	now := time.Now()

	cmd.Flags().StringVar(
		&opts.from,
		"from",
		now.Add(-15*time.Minute).Format(time.RFC3339),
		"Fetch logs from a specific time",
	)

	cmd.Flags().StringVar(
		&opts.to,
		"to",
		now.Format(time.RFC3339),
		"Fetch logs until a specific time",
	)

	cmd.Flags().Var(
		&opts.indexing,
		"index",
		"Indexing key-value pairs to filter logs (can be specified multiple times)",
	)

	return cmd
}
