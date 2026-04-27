import { IndicatorDot } from './styles'
import { PipelineTraceNodeIndicatorProps } from './types'

const PipelineTraceNodeIndicator = ({
  trace,
}: PipelineTraceNodeIndicatorProps) => {
  if (!trace) {
    return null
  }

  return <IndicatorDot hasError={trace.error !== null} />
}

export default PipelineTraceNodeIndicator
