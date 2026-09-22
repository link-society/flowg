import { useTranslation } from 'react-i18next'

import TextField from '@mui/material/TextField'

import BarChartIcon from '@mui/icons-material/BarChart'

import { Handle, NodeProps, Position } from '@xyflow/react'

import PipelineTraceNodeButton from '@/components/PipelineTraceNodeButton/component'
import PipelineTraceNodeIndicator from '@/components/PipelineTraceNodeIndicator/component'

import { NodeBody, NodeIcon, NodeRoot, ToolbarRow, handleStyle } from './styles'
import { PipelineNodeMetricData } from './types'

const PipelineNodeMetric = ({
  data,
  selected,
}: NodeProps<PipelineNodeMetricData>) => {
  const { t } = useTranslation()

  return (
    <>
      {selected && data.traces && (
        <ToolbarRow>
          <PipelineTraceNodeButton traces={data.traces} />
        </ToolbarRow>
      )}

      <Handle type="target" position={Position.Left} style={handleStyle} />
      <NodeRoot>
        <NodeIcon>
          <BarChartIcon />
        </NodeIcon>
        <NodeBody className="nodrag">
          <TextField
            label={t('components.pipelineNodeMetric.label')}
            type="text"
            value={data.name}
            slotProps={{
              input: {
                readOnly: true,
                sx: { fontFamily: 'monospace' },
              },
            }}
            variant="outlined"
          />
        </NodeBody>
      </NodeRoot>

      <PipelineTraceNodeIndicator
        status={
          data.traces
            ? data.traces.some((trace) => trace.error)
              ? 'error'
              : 'success'
            : null
        }
      />
    </>
  )
}

export default PipelineNodeMetric
