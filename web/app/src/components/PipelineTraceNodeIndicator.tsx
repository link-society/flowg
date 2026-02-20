import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

type PipelineTraceNodeIndicatorProps = {
  trace: NodeTrace | null
}

const PipelineTraceNodeIndicator = ({
  trace,
}: PipelineTraceNodeIndicatorProps) => (
  <>
    {trace && (
      <div
        style={{
          width: '18px',
          height: '18px',
          position: 'absolute',
          right: '-9px',
          top: '-9px',
          backgroundColor: trace.error === null ? '#20b834' : '#ff4444',
          borderRadius: '50%',
          boxShadow: '-2px 2px 2px #00000055',
        }}
      />
    )}
  </>
)

export default PipelineTraceNodeIndicator
