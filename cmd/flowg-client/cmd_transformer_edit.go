package main

import (
	"fmt"
	"os"

	"os/exec"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/shlex"
	"github.com/spf13/cobra"

	"link-society.com/flowg/api"

	"link-society.com/flowg/internal/utils/client"
)

func NewTransformerEditCommand() *cobra.Command {
	type options struct {
		name string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a transformer's code",
		Run: func(cmd *cobra.Command, args []string) {
			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := fmt.Sprintf("/api/v1/transformers/%s", opts.name)
			req, err := http.NewRequest(http.MethodGet, url, nil)
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

			var script string

			switch resp.StatusCode {
			case http.StatusOK:
				var data api.GetTransformerResponse
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
					exitCode = 1
					return
				}
				script = data.Script

			case http.StatusNotFound:
				script = ""

			default:
				fmt.Fprintf(os.Stderr, "ERROR: Received non-200 response: %s\n", resp.Status)
				exitCode = 1
				return
			}

			tmpf, err := os.CreateTemp("", "flowg-transformer-*.vrl")
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not create temporary file: %v\n", err)
				exitCode = 1
				return
			}
			defer os.Remove(tmpf.Name())
			defer tmpf.Close()

			if _, err := tmpf.WriteString(script); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not write to temporary file: %v\n", err)
				exitCode = 1
				return
			}
			tmpf.Close()

			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}

			editorArgv, err := shlex.Split(editor)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not parse editor command: %v\n", err)
				fmt.Fprintf(os.Stderr, "Command was: %s\n", editor)
				exitCode = 1
				return
			}
			editorArgv = append(editorArgv, tmpf.Name())

			cmdEditor := exec.Command(editorArgv[0], editorArgv[1:]...)
			cmdEditor.Stdin = os.Stdin
			cmdEditor.Stdout = os.Stdout
			cmdEditor.Stderr = os.Stderr
			if err := cmdEditor.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not run editor '%s': %v\n", editor, err)
				exitCode = 1
				return
			}

			content, err := os.ReadFile(tmpf.Name())
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not read temporary file: %v\n", err)
				exitCode = 1
				return
			}

			script = string(content)

			body := struct {
				Script string `json:"script"`
			}{
				Script: script,
			}

			payload, err := json.Marshal(body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not encode request body: %v\n", err)
				exitCode = 1
				return
			}

			req, err = http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not prepare request: %v\n", err)
				exitCode = 1
				return
			}

			resp, err = client.Do(req)
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
		},
	}

	cmd.Flags().StringVar(
		&opts.name,
		"name",
		"default",
		"Name of the transformer to delete",
	)

	return cmd
}
