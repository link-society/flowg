import { Node } from '@xyflow/react'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

export type PipelineNodeSourceData = Node<{
  type: string
  trace: NodeTrace | null
}>
