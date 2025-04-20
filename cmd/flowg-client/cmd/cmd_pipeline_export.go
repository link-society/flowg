package cmd

import (
	"fmt"
	"os"

	"html"
	"strings"

	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"

	"link-society.com/flowg/api"

	"link-society.com/flowg/internal/utils/client"
)

func NewPipelineExportCommand() *cobra.Command {
	type options struct {
		name   string
		format string
	}

	opts := &options{}

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			switch opts.format {
			case "mermaid":
			case "json":
			default:
				fmt.Fprintf(os.Stderr, "ERROR: Unsupported format: %s\n", opts.format)
				ExitCode = 1
				return
			}

			client := cmd.Context().Value(ApiClient).(*client.Client)
			url := fmt.Sprintf("/api/v1/pipelines/%s", opts.name)
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

			var data api.GetPipelineResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not decode response: %v\n", err)
				ExitCode = 1
				return
			}

			switch opts.format {
			case "json":
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				if err := encoder.Encode(data.Flow); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Could not encode response: %v\n", err)
					ExitCode = 1
					return
				}

			case "mermaid":
				fmt.Println("flowchart LR")
				for _, node := range data.Flow.Nodes {
					switch node.Type {
					case "source":
						fmt.Printf(
							"    %s[SOURCE: %s]\n",
							node.ID,
							strings.ReplaceAll(html.EscapeString(node.Data["type"]), "&#", "#"),
						)

					case "transform":
						fmt.Printf(
							"    %s[TRANSFORMER: %s]\n",
							node.ID,
							strings.ReplaceAll(html.EscapeString(node.Data["transformer"]), "&#", "#"),
						)

					case "switch":
						fmt.Printf(
							"    %s[SWITCH: %s]\n",
							node.ID,
							strings.ReplaceAll(html.EscapeString(node.Data["condition"]), "&#", "#"),
						)

					case "pipeline":
						fmt.Printf(
							"    %s[PIPELINE: %s]\n",
							node.ID,
							strings.ReplaceAll(html.EscapeString(node.Data["pipeline"]), "&#", "#"),
						)

					case "forwarder":
						fmt.Printf(
							"    %s[FORWARDER: %s]\n",
							node.ID,
							strings.ReplaceAll(html.EscapeString(node.Data["forwarder"]), "&#", "#"),
						)

					case "router":
						fmt.Printf(
							"    %s[STREAM: %s]\n",
							node.ID,
							strings.ReplaceAll(html.EscapeString(node.Data["stream"]), "&#", "#"),
						)
					}

					fmt.Printf("    class %s type-%s\n", node.ID, node.Type)
				}

				for _, edge := range data.Flow.Edges {
					fmt.Printf("    %s --> %s\n", edge.Source, edge.Target)
				}

				fmt.Printf("    classDef type-source fill:#ff6900,color:#fff,stroke:#c1380b,stroke-width:4px;\n")
				fmt.Printf("    classDef type-transform fill:#1547e6,color:#fff,stroke:#1d398d,stroke-width:4px;\n")
				fmt.Printf("    classDef type-switch fill:#e7000c,color:#fff,stroke:#c00109,stroke-width:4px;\n")
				fmt.Printf("    classDef type-pipeline fill:#efb300,color:#fff,stroke:#cf8701,stroke-width:4px;\n")
				fmt.Printf("    classDef type-forwarder fill:#008236,color:#fff,stroke:#0c552a,stroke-width:4px;\n")
				fmt.Printf("    classDef type-router fill:#8200dc,color:#fff,stroke:#57158d,stroke-width:4px;\n")
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.name,
		"name",
		"default",
		"Name of the pipeline to export",
	)

	cmd.Flags().StringVar(
		&opts.format,
		"format",
		"mermaid",
		"Format of the exported pipeline (mermaid, json)",
	)

	return cmd
}
