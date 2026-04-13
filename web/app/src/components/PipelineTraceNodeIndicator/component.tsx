import { IndicatorDot } from './styles'
import { PipelineTraceNodeIndicatorProps } from './types'

const PipelineTraceNodeIndicator = ({
  status,
}: PipelineTraceNodeIndicatorProps) => {
  if (status === null) {
    return null
  }

  return <IndicatorDot hasError={status === 'error'} />
}

export default PipelineTraceNodeIndicator
