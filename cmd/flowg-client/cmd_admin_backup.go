package main

import (
	"fmt"

	"os"
	"path"

	"io"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/utils/client"
)

func NewAdminBackupCommand() *cobra.Command {
	type options struct {
		dest string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup FlowG data",
		Run: func(cmd *cobra.Command, args []string) {
			if err := os.MkdirAll(opts.dest, os.ModePerm); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create backup directory: %v\n", err)
				exitCode = 1
				return
			}

			client := cmd.Context().Value(ApiClient).(*client.Client)
			if err := backup(client, "auth", opts.dest); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to backup auth database: %v\n", err)
				exitCode = 1
				return
			}

			if err := backup(client, "config", opts.dest); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to backup config database: %v\n", err)
				exitCode = 1
				return
			}

			if err := backup(client, "logs", opts.dest); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to backup logs database: %v\n", err)
				exitCode = 1
				return
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.dest,
		"dest",
		"./backup",
		"Destination for the backup",
	)

	return cmd
}

func backup(client *client.Client, dbType string, destDir string) error {
	url := fmt.Sprintf("/api/v1/backup/%s", dbType)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %r", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %s", resp.Status)
	}

	dest := path.Join(destDir, fmt.Sprintf("%s.db", dbType))
	outFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %r", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write backup file: %r", err)
	}

	return nil
}
