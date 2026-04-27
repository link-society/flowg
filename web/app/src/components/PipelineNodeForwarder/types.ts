import { Node } from '@xyflow/react'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

export type PipelineNodeForwarderData = Node<{
  forwarder: string
  trace: NodeTrace | null
}>
