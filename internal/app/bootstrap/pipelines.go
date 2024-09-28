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
				"nodes": [
					{
						"id": "__builtin__source_direct",
						"type": "source",
						"position": {"x": 210, "y": 195},
						"deletable": false,
						"data": {"type": "direct"},
						"measured": {"width": 136, "height": 38},
						"selected": true,
						"dragging": false
					},
					{
						"id": "__builtin__source_syslog",
						"type": "source",
						"position": {"x": 210, "y": 250},
						"deletable": false,
						"data": {"type": "syslog"},
						"measured": {"width": 136, "height": 38},
						"selected": true,
						"dragging": false
					},
					{
						"id": "node-1",
						"type": "router",
						"position": {"x": 405, "y": 195},
						"data": {"stream": "default"},
						"measured": {"width": 241,"height": 91},
						"selected": false,
						"dragging": false
					}
				],
				"edges": [
					{
						"id": "xy-edge____builtin__source_direct-node-1",
						"type": "smoothstep",
						"source": "__builtin__source_direct",
						"target": "node-1",
						"animated": true
					},
					{
						"id": "xy-edge____builtin__source_syslog-node-1",
						"type": "smoothstep",
						"source": "__builtin__source_syslog",
						"target": "node-1",
						"animated": true
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
