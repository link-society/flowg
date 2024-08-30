package bootstrap

import (
	"fmt"

	"link-society.com/flowg/internal/data/pipelines"
)

func DefaultPipeline(pipelinesManager *pipelines.Manager) error {
	pipelines, err := pipelinesManager.ListPipelines()
	if err != nil {
		return err
	}

	if len(pipelines) == 0 {
		err := pipelinesManager.SavePipelineFlow(
			"default",
			`{
				"nodes":[
					{
						"id": "__builtin__source",
						"type": "source",
						"position": {"x": 210, "y": 195},
						"deletable": false,
						"data": {},
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
				"edges":[
					{
						"id": "xy-edge____builtin__source-node-1",
						"type": "smoothstep",
						"source": "__builtin__source",
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
