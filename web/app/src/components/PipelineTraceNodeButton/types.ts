import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

export type PipelineTraceNodeButtonProps = {
  trace: NodeTrace
}

export type PipelineTraceNodeIndicatorProps = {
  trace: NodeTrace | null
}
