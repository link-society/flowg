import { Node } from '@xyflow/react'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

export type PipelineNodeMetricData = Node<{
  name: string
  traces: NodeTrace[] | null
}>
