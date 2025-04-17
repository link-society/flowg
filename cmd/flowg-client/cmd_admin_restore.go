package main

import (
	"fmt"
	"os"
	"path"

	"io"
	"mime/multipart"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/internal/utils/client"
)

func NewAdminRestoreCommand() *cobra.Command {
	type options struct {
		src string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore FlowG data",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := os.Stat(opts.src); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to find backup directory: %v\n", err)
				exitCode = 1
				return
			}

			client := cmd.Context().Value(ApiClient).(*client.Client)
			if err := restore(client, "auth", opts.src); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to restore auth database: %v\n", err)
				exitCode = 1
				return
			}

			if err := restore(client, "config", opts.src); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to restore config database: %v\n", err)
				exitCode = 1
				return
			}

			if err := restore(client, "logs", opts.src); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to restore logs database: %v\n", err)
				exitCode = 1
				return
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.src,
		"src",
		"./backup",
		"Source for the restore",
	)

	return cmd
}

func restore(client *client.Client, dbType string, srcDir string) error {
	srcPath := path.Join(srcDir, fmt.Sprintf("%s.db", dbType))

	body, writer := io.Pipe()

	url := fmt.Sprintf("/api/v1/restore/%s", dbType)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %w", err)
	}

	mw := multipart.NewWriter(writer)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	errC := make(chan error)
	go func() {
		defer close(errC)
		defer writer.Close()
		defer mw.Close()

		fw, err := mw.CreateFormFile("backup", srcPath)
		if err != nil {
			errC <- fmt.Errorf("failed to create form file: %w", err)
			return
		}

		file, err := os.Open(srcPath)
		if err != nil {
			errC <- fmt.Errorf("failed to open backup file: %w", err)
			return
		}
		defer file.Close()

		if n, err := io.Copy(fw, file); err != nil {
			errC <- fmt.Errorf("failed to copy backup file (%d bytes written): %w", n, err)
			return
		}

		if err := mw.Close(); err != nil {
			errC <- err
			return
		}
	}()

	resp, err := client.Do(req)
	merr := <-errC

	if err != nil || merr != nil {
		return fmt.Errorf("failed to send request: http error: %w, multipart error: %w", err, merr)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %s", resp.Status)
	}

	return nil
}
