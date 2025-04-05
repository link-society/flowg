package bootstrap

import (
	"context"
	"fmt"

	"link-society.com/flowg/internal/storage/config"
)

func DefaultPipeline(ctx context.Context, configStorage *config.Storage) error {
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
						"position": {
							"x": -45,
							"y": 0
						},
						"data": {
							"type": "direct"
						}
					},
					{
						"id": "__builtin__source_syslog",
						"type": "source",
						"position": {
							"x": -45,
							"y": 120
						},
						"data": {
							"type": "syslog"
						}
					},
					{
						"id": "node-37da64b3-6243-4620-b0f2-ba484a71b053",
						"type": "router",
						"position": {
							"x": 345,
							"y": 60
						},
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
