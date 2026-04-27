import { Node } from '@xyflow/react'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

export type PipelineNodeTransformerData = Node<{
  transformer: string
  trace: NodeTrace | null
}>
