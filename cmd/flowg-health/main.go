package main

import (
	"fmt"
	"net/http"
	"os"

	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/prometheus/procfs"
)

type options struct {
	pid int
}

var exitCode int = 0

func main() {
	opts := &options{}

	cmd := &cobra.Command{
		Use:   "flowg-health",
		Short: "Healthcheck for FlowG",
		Run: func(cmd *cobra.Command, args []string) {
			proc, err := procfs.NewProc(opts.pid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not retrieve process information: %v\n", err)
				exitCode = 1
				return
			}

			cmdline, err := proc.CmdLine()
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not retrieve process command line: %v\n", err)
				exitCode = 1
				return
			}

			var mgmtBindAddress string
			var mgmtTlsEnabled bool

			flags := pflag.NewFlagSet("flowg-server", pflag.ContinueOnError)
			flags.ParseErrorsWhitelist.UnknownFlags = true
			flags.StringVar(&mgmtBindAddress, "mgmt-bind", defaultMgmtBindAddress, "")
			flags.BoolVar(&mgmtTlsEnabled, "mgmt-tls", defaultMgmtTlsEnabled, "")

			if err := flags.Parse(cmdline); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not parse command line: %v\n", err)
				exitCode = 1
				return
			}

			host, port, err := net.SplitHostPort(mgmtBindAddress)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not parse host and port: %v\n", err)
				exitCode = 1
				return
			}

			if host == "" || host == "0.0.0.0" {
				host = "127.0.0.1"
			}

			var scheme string
			if mgmtTlsEnabled {
				scheme = "https"
			} else {
				scheme = "http"
			}

			url := fmt.Sprintf("%s://%s:%s/health", scheme, host, port)

			resp, err := http.Get(url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not send request: %v\n", err)
				exitCode = 1
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Fprintf(os.Stderr, "ERROR: Healthcheck failed with status code: %d\n", resp.StatusCode)
				exitCode = 1
				return
			}

			fmt.Println("OK")
		},
	}

	cmd.Flags().IntVar(
		&opts.pid,
		"pid",
		0,
		"PID of the FlowG process to check",
	)

	if err := cmd.Execute(); err != nil {
		exitCode = 1
	}

	os.Exit(exitCode)
}
