package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"link-society.com/flowg/cmd/flowg-client/utils"
)

// ExitCode is the process exit code; command handlers set it to signal failure.
var ExitCode int = 0

// NewRootCommand builds the root flowg-client command with its global flags and
// subcommands, constructing the API clients once before any subcommand runs.
func NewRootCommand() *cobra.Command {
	cobra.EnableTraverseRunHooks = true

	type globalOptions struct {
		apiUrl     string
		apiToken   string
		mgmtApiUrl string
	}

	opts := &globalOptions{}

	rootCmd := &cobra.Command{
		Use:   "flowg-client",
		Short: "API Client for FlowG",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ctx = context.WithValue(ctx, ApiClient, utils.NewClient(opts.apiUrl, opts.apiToken))
			ctx = context.WithValue(ctx, MgmtApiClient, utils.NewClient(opts.mgmtApiUrl, ""))
			cmd.SetContext(ctx)
		},
	}

	rootCmd.Flags().StringVar(
		&opts.apiUrl,
		"api",
		defaultApiUrl,
		"URL to FlowG HTTP API",
	)

	rootCmd.Flags().StringVar(
		&opts.apiToken,
		"api-token",
		defaultApiToken,
		"Authentication token for FlowG HTTP API",
	)

	rootCmd.Flags().StringVar(
		&opts.mgmtApiUrl,
		"mgmt-api",
		defaultMgmtApiUrl,
		"URL to FlowG Management HTTP API",
	)

	rootCmd.AddCommand(
		NewLoginCommand(),
		NewStreamCommand(),
		NewPipelineCommand(),
		NewTransformerCommand(),
		NewForwarderCommand(),
		NewAclCommand(),
		NewTokenCommand(),
		NewSystemConfigCommand(),
		NewAdminCommand(),
	)

	return rootCmd
}
