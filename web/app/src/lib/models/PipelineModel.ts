import { ReactFlowJsonObject } from '@xyflow/react'

type PipelineModel = {
  hasLayout: boolean
  nodes: ReactFlowJsonObject['nodes']
  edges: ReactFlowJsonObject['edges']
}

export default PipelineModel
