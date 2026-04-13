import { Node } from '@xyflow/react'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

export type PipelineNodeSwitchData = Node<{
  condition: string
  trace: NodeTrace | null
}>
