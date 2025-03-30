package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/api"

	"link-society.com/flowg/internal/utils/client"
)

func NewStreamListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List streams",
		Run: func(cmd *cobra.Command, args []string) {
			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := "/api/v1/streams"
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

			var data api.ListStreamsResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
				exitCode = 1
				return
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)

			fmt.Fprintf(w, "NAME\tRETENTION TIME\tRETENTION SIZE\tINDEXED FIELDS\n")

			for streamName, streamConfig := range data.Streams {
				var (
					retentionTime string
					retentionSize string
				)

				if streamConfig.RetentionTime == 0 {
					retentionTime = "disabled"
				} else {
					retentionTime = fmt.Sprintf("%ds", streamConfig.RetentionTime)
				}

				if streamConfig.RetentionSize == 0 {
					retentionSize = "disabled"
				} else {
					retentionSize = fmt.Sprintf("%dMB", streamConfig.RetentionSize)
				}

				fmt.Fprintf(
					w,
					"%s\t%s\t%s\t%s\n",
					streamName,
					retentionTime,
					retentionSize,
					strings.Join(streamConfig.IndexedFields, ","),
				)
			}

			if err := w.Flush(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not flush output: %v\n", err)
				exitCode = 1
				return
			}
		},
	}
}
