import { Node } from '@xyflow/react'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

export type PipelineNodeRouterData = Node<{
  stream: string
  trace: NodeTrace | null
}>
