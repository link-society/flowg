package bootstrap

import (
	"context"
	"fmt"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

func DefaultPipeline(ctx context.Context, configStorage storage.ConfigStorage) error {
	pipelines, err := configStorage.ListPipelines(ctx)
	if err != nil {
		return err
	}

	if len(pipelines) == 0 {
		err := configStorage.WriteRawPipeline(
			ctx,
			"default",
			`{
				"version": 2,
				"nodes": [
					{
						"id": "__builtin__source_direct",
						"type": "source",
						"data": {
							"type": "direct"
						}
					},
					{
						"id": "__builtin__source_syslog",
						"type": "source",
						"data": {
							"type": "syslog"
						}
					},
					{
						"id": "node-37da64b3-6243-4620-b0f2-ba484a71b053",
						"type": "router",
						"data": {
							"stream": "default"
						}
					}
				],
				"edges": [
					{
						"id": "xy-edge____builtin__source_syslog-node-37da64b3-6243-4620-b0f2-ba484a71b053",
						"source": "__builtin__source_syslog",
						"sourceHandle": "",
						"target": "node-37da64b3-6243-4620-b0f2-ba484a71b053"
					},
					{
						"id": "xy-edge____builtin__source_direct-node-37da64b3-6243-4620-b0f2-ba484a71b053",
						"source": "__builtin__source_direct",
						"sourceHandle": "",
						"target": "node-37da64b3-6243-4620-b0f2-ba484a71b053"
					}
				]
			}`,
		)
		if err != nil {
			return fmt.Errorf("failed to bootstrap default pipeline: %w", err)
		}
	}

	return nil
}
