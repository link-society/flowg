type PipelineTraceNodeIndicatorProps = {
  status: 'success' | 'error' | null
}

const PipelineTraceNodeIndicator = ({
  status,
}: PipelineTraceNodeIndicatorProps) => (
  <>
    {status !== null && (
      <div
        style={{
          width: '18px',
          height: '18px',
          position: 'absolute',
          right: '-9px',
          top: '-9px',
          backgroundColor: status === 'success' ? '#20b834' : '#ff4444',
          borderRadius: '50%',
          boxShadow: '-2px 2px 2px #00000055',
        }}
      />
    )}
  </>
)

export default PipelineTraceNodeIndicator
