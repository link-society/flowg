package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	client "link-society.com/flowg/cmd/flowg-client/cmd"
	health "link-society.com/flowg/cmd/flowg-health/cmd"
	server "link-society.com/flowg/cmd/flowg-server/cmd"
)

func main() {
	items := []struct {
		cmd        *cobra.Command
		destdir    string
		binaryName string
	}{
		{
			cmd:        client.NewRootCommand(),
			destdir:    "website/docs/cli/flowg-client",
			binaryName: "flowg-client",
		},
		{
			cmd:        health.NewRootCommand(),
			destdir:    "website/docs/cli/flowg-health",
			binaryName: "flowg-health",
		},
		{
			cmd:        server.NewRootCommand(),
			destdir:    "website/docs/cli/flowg-server",
			binaryName: "flowg-server",
		},
	}

	for _, item := range items {
		if err := genDoc(item.cmd, item.destdir, item.binaryName); err != nil {
			log.Fatal(err)
		}
	}
}

func genDoc(cmd *cobra.Command, destdir string, binaryName string) error {
	cmd.DisableAutoGenTag = true

	if err := os.MkdirAll(destdir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := doc.GenMarkdownTreeCustom(cmd, destdir, docusaurusHeader, docusaurusLink); err != nil {
		return fmt.Errorf("failed to generate documentation for `%s`: %w", binaryName, err)
	}

	return nil
}

func docusaurusHeader(filename string) string {
	name := filepath.Base(filename)
	base := strings.TrimSuffix(name, filepath.Ext(name))
	title := strings.ReplaceAll(base, "_", " ")

	return fmt.Sprintf(headerTemplate, title)
}

func docusaurusLink(filename string) string {
	return filename
}

const headerTemplate = `---
title: %s
hide_title: true
---
`
