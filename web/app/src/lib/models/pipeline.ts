import { ReactFlowJsonObject } from '@xyflow/react'

export type PipelineModel = {
  nodes: ReactFlowJsonObject['nodes'],
  edges: ReactFlowJsonObject['edges'],
}
