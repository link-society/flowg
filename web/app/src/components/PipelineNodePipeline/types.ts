import { Node } from '@xyflow/react'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

export type PipelineNodePipelineData = Node<{
  pipeline: string
  trace: NodeTrace | null
}>
