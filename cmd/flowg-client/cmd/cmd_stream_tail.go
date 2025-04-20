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
	"link-society.com/flowg/internal/utils/client/log"
	"link-society.com/flowg/internal/utils/timex"
)

func NewStreamTailCommand() *cobra.Command {
	type options struct {
		name   string
		filter string
		period string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "tail",
		Short: "Fetch logs until now",
		Run: func(cmd *cobra.Command, args []string) {
			period, err := timex.ParseDuration(opts.period)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not parse period: %v\n", err)
				ExitCode = 1
				return
			}

			to := time.Now()
			from := to.Add(-period)

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

			queryset.Set("from", from.Format(time.RFC3339))
			queryset.Set("to", to.Format(time.RFC3339))

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

	cmd.Flags().StringVar(
		&opts.period,
		"period",
		"15m",
		"Timespan to fetch logs for",
	)

	return cmd
}
