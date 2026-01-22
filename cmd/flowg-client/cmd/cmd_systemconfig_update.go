package cmd

import (
	"fmt"
	"os"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/api"

	"link-society.com/flowg/internal/utils/client"
)

func NewSystemConfigUpdateCommand() *cobra.Command {
	type options struct {
		SyslogAllowedOrigins []string
		SyslogAllowAll       bool
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update system configuration",
		Run: func(cmd *cobra.Command, args []string) {
			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := "/api/v1/system-configuration"
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
				fmt.Fprintf(os.Stderr, "ERROR: Received non-200 response: %s\n", resp.Status)
				ExitCode = 1
				return
			}

			var data api.GetSystemConfigurationResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
				ExitCode = 1
				return
			}

			if cmd.Flags().Changed("syslog-allowed-origins") {
				data.Configuration.SyslogAllowedOrigins = opts.SyslogAllowedOrigins
			}

			if cmd.Flags().Changed("syslog-allow-all") && opts.SyslogAllowAll {
				data.Configuration.SyslogAllowedOrigins = []string{}
			}

			payload, err := json.Marshal(data.Configuration)
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
				fmt.Fprintf(os.Stderr, "ERROR: Failed to update system configuration: %s\n", resp.Status)
				ExitCode = 1
				return
			}
		},
	}

	cmd.Flags().StringArrayVar(
		&opts.SyslogAllowedOrigins,
		"syslog-allowed-origin",
		[]string{},
		"List of allowed origins for syslog (empty list means all origins are allowed)",
	)

	cmd.Flags().BoolVar(
		&opts.SyslogAllowAll,
		"syslog-allow-all",
		false,
		"Allow all origins for syslog (overrides --syslog-allowed-origin)",
	)

	return cmd
}
