package cmd

import (
	"fmt"
	"os"

	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/api/schemas"
	"link-society.com/flowg/cmd/flowg-client/utils"
)

// NewSystemConfigShowCommand builds the "show" command, which shows the system configuration.
func NewSystemConfigShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show system configuration",
		Run: func(cmd *cobra.Command, args []string) {
			client := cmd.Context().Value(ApiClient).(*utils.Client)
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

			var data schemas.GetSystemConfigurationResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
				ExitCode = 1
				return
			}

			fmt.Println("System Configuration:")

			if len(data.Configuration.SyslogAllowedOrigins) == 0 {
				fmt.Println("  Syslog Allowed Origins: ALL")
			} else {
				fmt.Println("  Syslog Allowed Origins:")
				for _, origin := range data.Configuration.SyslogAllowedOrigins {
					fmt.Printf("    - %s\n", origin)
				}
			}
		},
	}
}
