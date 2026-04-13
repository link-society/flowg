import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

export type PipelineTraceNodeButtonProps = {
  traces: NodeTrace[]
}

export type PipelineTraceNodeIndicatorProps = {
  traces: NodeTrace[] | null
}
