import { Node } from '@xyflow/react'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

export type PipelineNodeSourceData = Node<{
  type: string
  traces: NodeTrace[] | null
}>
